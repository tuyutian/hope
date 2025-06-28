package users

import (
	"context"
	"errors"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"backend/internal/domain/entity/users"
	"backend/internal/domain/repo"
	appRepo "backend/internal/domain/repo/apps"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/pkg/response/message"
)

type UserService struct {
	appRepo     appRepo.AppRepository
	userRepo    userRepo.UserRepository
	cacheRepo   repo.CacheRepository
	userCache   userRepo.UserCacheRepository
	shopifyRepo shopifyRepo.ShopifyRepository
	appAuthRepo userRepo.AppAuthRepository
}

func NewUserService(repos *providers.Repositories) *UserService {
	return &UserService{userRepo: repos.UserRepo}
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

func (u *UserService) AuthFromSession(ctx context.Context, token string) (*users.User, error) {
	appID := ctx.Value(ctxkeys.AppID).(string)
	claims := u.GetClaims(ctx)
	sessionToken, err := u.shopifyRepo.RequestOfflineSessionToken(ctx, token)
	if err != nil {
		return nil, err
	}
	app := ctx.Value(ctxkeys.ShopifyApp).(goshopify.App)
	shopName, err := u.shopifyRepo.GetShopName(ctx, claims.Dest)
	if err != nil {
		return nil, err
	}
	client, err := app.NewClient(shopName, sessionToken.Token, nil)
	shop, err := client.Shop.Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	// 恢复用户数据
	user, _ := u.userRepo.GetByShop(ctx, appID, claims.Dest)
	user.AppId = appID
	user.AccessToken = sessionToken.Token
	user.Shop = claims.Dest
	user.IsDel = 0
	if shop != nil {
		user.City = shop.City
		user.CurrencyCode = shop.Country
		user.CountryName = shop.CountryName
		user.CountryCode = shop.CountryCode
		user.Email = shop.Email
		user.Phone = shop.Phone
		user.Timezone = shop.Timezone
		user.PlanDisplayName = shop.PlanDisplayName
		user.Name = shop.Name
		user.CurrencyCode = shop.Currency
		user.MoneyFormat = u.shopifyRepo.ExtractCurrencySymbol(shop.MoneyFormat)
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
		user, err = u.userRepo.Get(ctx, id, "id", "username", "email")
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
