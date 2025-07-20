package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/entity/users"
	userEntity "backend/internal/domain/entity/users"
	"backend/internal/domain/repo"
	appRepo "backend/internal/domain/repo/apps"
	"backend/internal/domain/repo/jobs"
	jwtRepo "backend/internal/domain/repo/jwtauth"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

type UserService struct {
	appRepo            appRepo.AppRepository
	userRepo           userRepo.UserRepository
	cacheRepo          repo.CacheRepository
	userCache          userRepo.UserCacheRepository
	userSettingRepo    userRepo.UserSettingRepository
	shopifyRepo        shopifyRepo.ShopifyRepository
	appAuthRepo        appRepo.AppAuthRepository
	shopGraphqlRepo    shopifyRepo.ShopGraphqlRepository
	productGraphqlRepo shopifyRepo.ProductGraphqlRepository
	asynqRepo          jobs.AsynqRepository
	jwtRepo            jwtRepo.JWTRepository
	subscriptionRepo   userRepo.UserSubscriptionRepository
}

func NewUserService(repos *providers.Repositories) *UserService {
	return &UserService{
		userRepo:           repos.UserRepo,
		appRepo:            repos.AppRepo,
		cacheRepo:          repos.CacheRepo,
		userCache:          repos.UserCacheRepo,
		shopifyRepo:        repos.ShopifyRepo,
		appAuthRepo:        repos.AppAuthRepo,
		productGraphqlRepo: repos.ProductGraphqlRepo,
		shopGraphqlRepo:    repos.ShopGraphqlRepo,
		asynqRepo:          repos.AsyncRepo,
		jwtRepo:            repos.JwtRepo,
		userSettingRepo:    repos.UserSettingRepo,
		subscriptionRepo:   repos.UserSubscriptionRepo,
	}
}

func (u *UserService) GetLoginUserFromID(ctx context.Context, id int64) (*users.User, error) {
	user, err := u.userRepo.Get(ctx, id)
	return user, err
}

func (u *UserService) GetLoginUserFromShop(ctx context.Context, shop string) (*users.User, error) {
	appId := ctx.Value(ctxkeys.AppID).(string)
	user, err := u.userRepo.GetActiveUserByShop(ctx, appId, shop)
	return user, err
}

func (u *UserService) GetLoginAdminFromID(ctx context.Context, id int64) (interface{}, error) {
	return nil, nil
}

// GetClaims 从 context 中获取 jwt.BizClaims
func (u *UserService) GetClaims(ctx context.Context) *jwt.BizClaims {
	claims, _ := ctx.Value(ctxkeys.BizClaims).(*jwt.BizClaims)
	return claims
}

// GetShopifyClient 从 context 中获取 client
func (u *UserService) GetShopifyClient(ctx context.Context) *shopify_graphql.GraphqlClient {
	client, _ := ctx.Value(ctxkeys.ShopifyGraphqlClient).(*shopify_graphql.GraphqlClient)
	return client
}

func (u *UserService) AuthFromSession(ctx context.Context, sessionToken *shopifyEntity.Token, claims *jwt.BizClaims) (*users.User, error) {
	appID := ctx.Value(ctxkeys.AppID).(string)
	shopName, err := utils.GetShopName(claims.Dest)
	if err != nil {
		return nil, err
	}
	client := shopify_graphql.NewGraphqlClient(shopName, sessionToken.Token)
	u.shopGraphqlRepo.WithClient(client)
	shop, err := u.shopGraphqlRepo.GetShopInfo(ctx)
	if err != nil {
		logger.Error(ctx, "shopify_graphql_repo.GetShopInfo", zap.Error(err))
	}
	// 恢复用户数据
	user, _ := u.userRepo.GetByShop(ctx, appID, claims.Dest)
	if user == nil {
		user = &users.User{}
	}
	user.AppId = appID
	user.AccessToken = sessionToken.Token
	user.Shop = claims.Dest
	user.IsDel = 0
	if shop != nil {
		user.City = shop.BillingAddress.City
		user.CountryName = shop.BillingAddress.Country
		user.CountryCode = shop.BillingAddress.CountryCodeV2
		user.Email = shop.Email
		user.Phone = shop.BillingAddress.Phone
		user.Timezone = shop.TimezoneOffsetMinutes
		user.PlanDisplayName = shop.Plan.DisplayName
		user.Name = shop.Name
		user.CurrencyCode = shop.CurrencyCode
		user.MoneyFormat = u.shopifyRepo.ExtractCurrencySymbol(shop.CurrencyFormats.MoneyFormat)
	}
	if user.ID > 0 {
		err := u.userRepo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	} else {
		id, err := u.userRepo.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}
		user.ID = id
	}
	// 更新 app auth记录
	appAuth, err := u.appAuthRepo.GetByUserAndApp(ctx, user.ID, appID)
	if err == nil {
		appAuth.Scopes = sessionToken.Scope
		appAuth.Shop = claims.Dest
		appAuth.Status = 1
		if appAuth.Id > 0 {
			_, _ = u.appAuthRepo.Create(ctx, appAuth)
		} else {
			_ = u.appAuthRepo.Update(ctx, appAuth)
		}
	}
	// 注册必要地初始化数据

	return user, err
}

// getUser 先从缓存中获取 entity.UserInfo，如果不存在，则从数据库获取
func (u *UserService) getUser(ctx context.Context, id int64) (*users.User, error) {
	// 从缓存中获取
	user, err := u.userCache.Get(ctx, id)
	// key 不存在, 查询数据库并更新缓存
	if errors.Is(err, redis.Nil) {
		user, err = u.userRepo.Get(ctx, id, "id", "shop", "app_id", "plans", "email", "name", "timezone", "access_token", "money_format", "is_del")
		if err != nil {
			return nil, err
		}
		// 缓存无效用户，防止 redis 穿透
		if user.ID == 0 {
			invalidUser := &users.User{ID: -1}
			err = u.userCache.Set(ctx, id, invalidUser, 30*time.Second)
			if err != nil {
				logger.Warn(ctx, "failed to execute redis.set", zap.Error(err))
			}
			return nil, message.ErrInvalidAccount
		}
		// 更新缓存
		if err := u.userCache.Set(ctx, id, user, 10*time.Minute); err != nil {
			logger.Warn(ctx, "failed to execute redis.set", zap.Error(err))
		}
		return user, nil
	}
	// other error
	if err != nil {
		return nil, err
	}
	// key 存在，但是用户不存在
	if user.ID < 1 {
		return nil, message.ErrInvalidAccount
	}
	return user, nil
}

func (u *UserService) Uninstall(ctx context.Context, appId string, shop string) error {
	user, err := u.userRepo.GetActiveUserByShop(ctx, appId, shop)

	// 如果未安装
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	err = u.userRepo.UpdateIsDel(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "uninstall db异常", "Err:", err.Error())
		return err
	}

	// 卸载 关闭购物车 清空状态 删除shopify产品
	_, err = u.asynqRepo.DelProductTask(ctx, user.ID, 0, 1)

	return nil
}

func (u *UserService) UpdateUserStep(ctx context.Context, step userEntity.UpdateStep) error {
	claims := u.GetClaims(ctx)
	var steps map[string]bool
	userID := claims.UserID
	stepSettingStr, err := u.userSettingRepo.Get(ctx, userID, userEntity.DashboardGuideStep)

	if stepSettingStr != "" {
		err = json.Unmarshal([]byte(stepSettingStr), &steps)
	} else {
		steps = userEntity.DefaultDashboardGuideStep
	}
	steps[step.Name] = step.Open
	// 4. 再把 steps 转成字符串，保存回数据库
	updatedSteps, err := json.Marshal(steps)
	if err != nil {
		logger.Error(ctx, "update-step json(2)异常", "Err:", err.Error())
		return err
	}
	return u.userSettingRepo.Set(ctx, userID, userEntity.DashboardGuideStep, string(updatedSteps))
}

type CollectionOption struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}
type UserConfigResponse struct {
	MoneySymbol  string `json:"money_symbol"`
	HasSubscribe bool   `json:"has_subscribe"`
}
type Collection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (u *UserService) GetUserConf(ctx context.Context, userID int64) (*UserConfigResponse, error) {
	user, err := u.getUser(ctx, userID)

	if err != nil {
		logger.Error(ctx, "get-user-info db异常", "Err:", err.Error())
		return nil, err
	}

	subscribe, _ := u.subscriptionRepo.GetActiveSubscription(ctx, userID)

	return &UserConfigResponse{
		MoneySymbol:  user.MoneyFormat,
		HasSubscribe: subscribe != nil,
	}, nil
}

func (u *UserService) GetSessionData(ctx context.Context, userID int64) (*userEntity.SessionData, error) {
	user, err := u.getUser(ctx, userID)

	if err != nil {
		logger.Error(ctx, "get-user-info db异常", "Err:", err.Error())
		return nil, err
	}
	var steps map[string]bool
	stepSettingStr, _ := u.userSettingRepo.Get(ctx, userID, userEntity.DashboardGuideStep)
	if stepSettingStr != "" {
		err = json.Unmarshal([]byte(stepSettingStr), &steps)
	} else {
		steps = userEntity.DefaultDashboardGuideStep
	}
	guideHide, _ := u.userSettingRepo.Get(ctx, userID, userEntity.DashboardGuideHide)
	return &userEntity.SessionData{
		Shop:      user.Shop,
		GuideStep: steps,
		GuideShow: len(guideHide) == 0 || guideHide == "0",
	}, nil
}

func (u *UserService) SyncShopifyUserInfo(ctx context.Context, shop string, planDisplayName string) error {
	user, err := u.userRepo.FirstName(ctx, shop)

	if err != nil {
		logger.Error(ctx, "sync-user-info db异常", "Err:", err.Error())
		return err
	}

	if user == nil || user.IsDel != 1 {
		return nil
	}

	// 检测是不是关店了
	if planDisplayName == "Frozen" || planDisplayName == "Cancelled" {
		err = u.userRepo.UpdateIsClose(ctx, user.ID, planDisplayName)
		if err != nil {
			return err
		}

		// 卸载 关闭购物车 清空状态 删除shopify产品
		_, err := u.asynqRepo.DelProductTask(ctx, user.ID, 0, 1)

		if err != nil {
			logger.Error(ctx, "关店 del_product_queue 推送队列失败:", err.Error())
			return err
		}

		return nil
	}

	u.shopGraphqlRepo.WithClient(u.GetShopifyClient(ctx))

	// 拿到Token 需要去获取用户基本信息
	shopInfo, err := u.shopGraphqlRepo.GetShopInfo(ctx) // 通过 client 调用方法

	if err != nil {
		return fmt.Errorf("获取店铺信息异常: %w", err)
	}

	var userModel users.User
	userModel.ID = user.ID
	userModel.City = shopInfo.BillingAddress.City
	userModel.CountryCode = shopInfo.BillingAddress.CountryCodeV2
	userModel.CountryName = shopInfo.BillingAddress.Country
	userModel.CurrencyCode = shopInfo.CurrencyCode
	userModel.MoneyFormat = u.shopifyRepo.ExtractCurrencySymbol(shopInfo.CurrencyFormats.MoneyFormat)
	userModel.PlanDisplayName = shopInfo.Plan.DisplayName

	_ = u.userRepo.Update(ctx, &userModel)

	return nil
}

func (u *UserService) GenerateTestToken(ctx context.Context, id int64) string {
	user, _ := u.getUser(ctx, id)
	claims := jwt.BizClaims{
		Dest:    user.Shop,
		UserID:  user.ID,
		AdminID: 0,
		Jti:     utils.Uuid(),
		Sid:     "",
		Sig:     "",
	}
	t, _, _ := u.jwtRepo.GenerateToken(ctx, claims)
	return t
}

func (u *UserService) UpdateUserSetting(ctx context.Context, setting userEntity.UpdateSetting) error {
	claims := u.GetClaims(ctx)
	userID := claims.UserID
	return u.userSettingRepo.Set(ctx, userID, setting.Name, setting.Value)
}
