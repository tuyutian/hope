package users

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	appEntity "backend/internal/domain/entity/apps"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	userEntity "backend/internal/domain/entity/users"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type SubscriptionService struct {
	userSubscriptionRepo    users.UserSubscriptionRepository
	subscriptionGraphqlRepo shopifyRepo.SubscriptionGraphqlRepository
	usageChargeGraphqlRepo  shopifyRepo.UsageChargeGraphqlRepository
	shopifyRepo             shopifyRepo.ShopifyRepository
}

func NewSubscriptionService(
	repos *providers.Repositories,
) *SubscriptionService {
	return &SubscriptionService{
		userSubscriptionRepo:    repos.UserSubscriptionRepo,
		subscriptionGraphqlRepo: repos.SubscriptionGraphqlRepo,
		usageChargeGraphqlRepo:  repos.UsageChargeGraphqlRepo,
		shopifyRepo:             repos.ShopifyRepo,
	}
}

// CreateUsageSubscription 创建用量订阅
func (s *SubscriptionService) CreateUsageSubscription(
	ctx context.Context,
	planName string,
	cappedAmount decimal.Decimal,
	currency string,
	terms string,
	isTest bool,
) (*userEntity.UserSubscription, string, error) {
	claims := ctx.Value(ctxkeys.BizClaims).(*jwt.BizClaims)
	appData := ctx.Value(ctxkeys.AppData).(appEntity.AppData)
	returnUrl := s.shopifyRepo.GetReturnUrl(appData.AppID, claims.UserID)
	// 1. 构建订阅输入
	input := shopifyEntity.AppSubscriptionCreateInput{
		Name: planName,
		Test: isTest,
		LineItems: []shopifyEntity.AppSubscriptionLineItemInput{
			{
				Plan: shopifyEntity.AppPlanInput{
					AppUsagePricingDetails: &shopifyEntity.AppUsagePricingDetailsInput{
						CappedAmount: shopifyEntity.MoneyInput{
							Amount:       cappedAmount,
							CurrencyCode: currency,
						},
						Terms: terms,
					},
					AppRecurringPricingDetails: nil,
				},
			},
		},
		ReturnURL: returnUrl,
	}

	// 2. 创建 Shopify 订阅
	s.subscriptionGraphqlRepo.WithClient(ctx.Value(ctxkeys.ShopifyGraphqlClient).(*shopify_graphql.GraphqlClient))
	subscription, confirmationURL, err := s.subscriptionGraphqlRepo.CreateSubscription(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Shopify subscription: %v", err)
	}
	// 3. 保存到本地数据库
	userSubscription := &userEntity.UserSubscription{
		UserID:                 claims.UserID,
		ShopDomain:             claims.Dest,
		ChargeID:               utils.GetIdFromShopifyGraphqlId(subscription.ID),
		SubscriptionName:       subscription.Name,
		SubscriptionStatus:     subscription.Status,
		SubscriptionLineItemID: subscription.LineItems[0].ID,
		PricingType:            userEntity.PricingTypeRecurring,
		CappedAmount:           &cappedAmount,
		Currency:               currency,
		BalanceUsed:            &decimal.Zero,
		Price:                  &decimal.Zero,
		Terms:                  terms,
		CurrentPeriodStart:     time.Now().Unix(),
		CurrentPeriodEnd:       utils.ParseShopifyTime(subscription.CurrentPeriodEnd),
		TrialDays:              subscription.TrialDays,
		TestSubscription:       subscription.Test,
		LastSyncTime:           time.Now().Unix(),
	}

	err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, userSubscription)
	if err != nil {
		return nil, "", fmt.Errorf("failed to save subscription to database: %v", err)
	}

	return userSubscription, confirmationURL, nil
}

// CreateRecurringSubscription 创建循环订阅
func (s *SubscriptionService) CreateRecurringSubscription(
	ctx context.Context,
	userID int64,
	shopDomain string,
	planName string,
	price decimal.Decimal,
	currency string,
	interval string, // EVERY_30_DAYS, ANNUAL
	returnURL string,
	isTest bool,
) (*userEntity.UserSubscription, string, error) {

	// 1. 构建订阅输入
	input := shopifyEntity.AppSubscriptionCreateInput{
		Name: planName,
		Test: isTest,
		LineItems: []shopifyEntity.AppSubscriptionLineItemInput{
			{
				Plan: shopifyEntity.AppPlanInput{
					AppRecurringPricingDetails: &shopifyEntity.AppRecurringPricingDetailsInput{
						Price: shopifyEntity.MoneyInput{
							Amount:       price,
							CurrencyCode: currency,
						},
						Interval: interval,
					},
				},
			},
		},
		ReturnURL: returnURL,
	}

	// 2. 创建 Shopify 订阅
	s.subscriptionGraphqlRepo.WithClient(ctx.Value(ctxkeys.ShopifyGraphqlClient).(*shopify_graphql.GraphqlClient))
	subscription, confirmationURL, err := s.subscriptionGraphqlRepo.CreateSubscription(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Shopify subscription: %v", err)
	}

	// 3. 保存到本地数据库
	userSubscription := &userEntity.UserSubscription{
		UserID:                 userID,
		ShopDomain:             shopDomain,
		ChargeID:               utils.GetIdFromShopifyGraphqlId(subscription.ID),
		SubscriptionName:       subscription.Name,
		SubscriptionStatus:     subscription.Status,
		SubscriptionLineItemID: subscription.LineItems[0].ID,
		PricingType:            userEntity.PricingTypeRecurring,
		Price:                  &price, // 对于循环订阅，这里存储价格
		Currency:               currency,
		CappedAmount:           &decimal.Zero,
		BalanceUsed:            &decimal.Zero,
		Terms:                  fmt.Sprintf("Recurring charge - %s", interval),
		CurrentPeriodStart:     time.Now().Unix(),
		CurrentPeriodEnd:       utils.ParseShopifyTime(subscription.CurrentPeriodEnd),
		TrialDays:              subscription.TrialDays,
		TestSubscription:       subscription.Test,
		LastSyncTime:           time.Now().Unix(),
	}

	err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, userSubscription)
	if err != nil {
		return nil, "", fmt.Errorf("failed to save subscription to database: %v", err)
	}

	return userSubscription, confirmationURL, nil
}

// SyncSubscriptionStatus 同步订阅状态
func (s *SubscriptionService) SyncSubscriptionStatus(ctx context.Context, user *userEntity.User) error {
	// 1. 从 Shopify 获取当前订阅
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	s.subscriptionGraphqlRepo.WithClient(client)
	currentSubscription, err := s.subscriptionGraphqlRepo.GetCurrentSubscription(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current subscription from Shopify: %v", err)
	}

	// 解析当前订阅的 chargeID
	currentChargeID := utils.GetIdFromShopifyGraphqlId(currentSubscription.ID)

	// 2. 先更新或创建当前订阅
	userSubscription, err := s.userSubscriptionRepo.GetSubscriptionByChargeID(ctx, currentChargeID)
	if err != nil {
		return fmt.Errorf("failed to get user subscription from database: %v", err)
	}

	var newSubscription *userEntity.UserSubscription

	if userSubscription != nil {
		// 更新现有订阅
		userSubscription.SubscriptionStatus = currentSubscription.Status
		userSubscription.CurrentPeriodEnd = utils.ParseShopifyTime(currentSubscription.CurrentPeriodEnd)
		userSubscription.LastSyncTime = time.Now().Unix()
		userSubscription.UpdateTime = time.Now().Unix()

		err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, userSubscription)
		if err != nil {
			return fmt.Errorf("failed to update user subscription: %v", err)
		}
		newSubscription = userSubscription
	} else {
		// 创建新地订阅记录
		if len(currentSubscription.LineItems) == 0 {
			return fmt.Errorf("no line items found in current subscription")
		}

		// 解析订阅 ID
		lineItemID := currentSubscription.LineItems[0].ID

		// 确定定价类型和相关信息
		var pricingType string
		var cappedAmount decimal.Decimal
		var terms string
		var balanceUsed decimal.Decimal
		var currency = "USD" // 默认货币

		if currentSubscription.IsUsageSubscription() {
			pricingType = userEntity.PricingTypeRecurring
			usagePricing, _ := currentSubscription.GetUsagePricing()
			cappedAmount, _ = decimal.NewFromString(usagePricing.CappedAmount.Amount)
			terms = usagePricing.Terms
			balanceUsed, _ = decimal.NewFromString(usagePricing.BalanceUsed.Amount)
			currency = usagePricing.CappedAmount.CurrencyCode
		} else if currentSubscription.IsRecurringSubscription() {
			pricingType = userEntity.PricingTypeRecurring
			recurringPricing, _ := currentSubscription.GetRecurringPricing()
			cappedAmount, _ = decimal.NewFromString(recurringPricing.Price.Amount)
			terms = fmt.Sprintf("Recurring subscription with %s interval", recurringPricing.Interval)
			balanceUsed = decimal.Zero
			currency = recurringPricing.Price.CurrencyCode
		}

		// 创建新的订阅记录
		newSubscription = &userEntity.UserSubscription{
			UserID:                 user.ID,
			ShopDomain:             user.Shop,
			ChargeID:               currentChargeID,
			SubscriptionName:       currentSubscription.Name,
			SubscriptionStatus:     currentSubscription.Status,
			SubscriptionLineItemID: lineItemID,
			PricingType:            pricingType,
			CappedAmount:           &cappedAmount,
			Currency:               currency,
			BalanceUsed:            &balanceUsed,
			Terms:                  terms,
			CurrentPeriodStart:     utils.ParseShopifyTime(currentSubscription.CreatedAt),
			CurrentPeriodEnd:       utils.ParseShopifyTime(currentSubscription.CurrentPeriodEnd),
			TrialDays:              currentSubscription.TrialDays,
			TestSubscription:       currentSubscription.Test,
			LastSyncTime:           time.Now().Unix(),
			CreateTime:             time.Now().Unix(),
			UpdateTime:             time.Now().Unix(),
		}

		err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, newSubscription)
		if err != nil {
			return fmt.Errorf("failed to create new user subscription: %v", err)
		}
	}

	// 3. 只有在当前订阅成功更新/创建后，才取消其他活跃订阅
	// 并且只有当当前订阅状态是 ACTIVE 时才执行这个操作
	if newSubscription.SubscriptionStatus == userEntity.SubscriptionStatusActive {
		err = s.userSubscriptionRepo.CancelActiveSubscriptionsExcept(ctx, user.ID, currentChargeID)
		if err != nil {
			// 这里可以记录警告日志，但不返回错误，因为主要任务已经完成
			fmt.Printf("Warning: failed to cancel other active subscriptions for user %d: %v\n", user.ID, err)
		}
	}

	return nil
}

func (s *SubscriptionService) HandleApproachingCappedAmount(ctx context.Context, chargeID int64) error {
	logger.Warn(ctx, "usage capped amount approaching: ", chargeID)
	return nil
}

func (s *SubscriptionService) VerifyPayment(ctx context.Context, user *userEntity.User, chargeID int64) (*userEntity.UserSubscription, error) {
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)

	s.subscriptionGraphqlRepo.WithClient(client)
	subscription, err := s.subscriptionGraphqlRepo.GetRecurrentChargeByID(ctx, chargeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription from Shopify: %v", err)
	}
	if subscription == nil {
		return nil, fmt.Errorf("fetch subscription failed")
	}
	usagePrice, _ := subscription.GetUsagePricing()
	CappedAmount, _ := decimal.NewFromString(usagePrice.CappedAmount.Amount)
	BalanceUsed, _ := decimal.NewFromString(usagePrice.BalanceUsed.Amount)
	Terms := usagePrice.Terms
	fmt.Println(subscription.LineItems[0].ID)
	// 3. 保存到本地数据库
	userSubscription := &userEntity.UserSubscription{
		UserID:                 user.ID,
		ShopDomain:             user.Shop,
		ChargeID:               utils.GetIdFromShopifyGraphqlId(subscription.ID),
		SubscriptionName:       subscription.Name,
		SubscriptionStatus:     subscription.Status,
		SubscriptionLineItemID: subscription.LineItems[0].ID,
		PricingType:            userEntity.PricingTypeRecurring,
		CappedAmount:           &CappedAmount, // 对于循环订阅，这里存储价格
		Currency:               usagePrice.CappedAmount.CurrencyCode,
		BalanceUsed:            &BalanceUsed,
		Terms:                  Terms,
		CurrentPeriodStart:     time.Now().Unix(),
		CurrentPeriodEnd:       utils.ParseShopifyTime(subscription.CurrentPeriodEnd),
		TrialDays:              subscription.TrialDays,
		TestSubscription:       subscription.Test,
		LastSyncTime:           time.Now().Unix(),
	}
	err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, userSubscription)
	if err != nil {
		return nil, fmt.Errorf("failed to save subscription to database: %v", err)
	}
	return userSubscription, nil
}

func (s *SubscriptionService) UpdateSubscriptionStatus(ctx context.Context, user *userEntity.User, chargeID int64, status string) error {
	err := s.userSubscriptionRepo.UpdateSubscriptionStatus(ctx, chargeID, status)
	if err != nil {
		return err
	}
	if user != nil && status == userEntity.SubscriptionStatusActive {
		err := s.SyncSubscriptionStatus(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SubscriptionService) CreateUsageCharge(ctx context.Context, user *userEntity.User, lineItemId string, amount decimal.Decimal, currency string) error {
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	s.usageChargeGraphqlRepo.WithClient(client)
	usageRecordID, err := s.usageChargeGraphqlRepo.CreateUsageCharge(ctx, lineItemId, amount, currency)
	if err != nil {
		logger.Error(ctx, "failed to create usage charge: ", err)
		return err
	}
	if usageRecordID == "" {
		utils.CallWilding(fmt.Sprintf("failed to create usage charge with user: %d amount:%v ", user.ID, amount))
		return fmt.Errorf("failed to create usage charge")
	}

	return nil
}
