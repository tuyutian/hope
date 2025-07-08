package users

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	billingEntity "backend/internal/domain/entity/billings"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/billings"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/providers"
)

type SubscriptionService struct {
	userSubscriptionRepo    billings.UserSubscriptionRepository
	subscriptionGraphqlRepo shopifyRepo.SubscriptionGraphqlRepository
	usageChargeGraphqlRepo  shopifyRepo.UsageChargeGraphqlRepository
}

func NewSubscriptionService(
	repos *providers.Repositories,
) *SubscriptionService {
	return &SubscriptionService{
		userSubscriptionRepo:    repos.UserSubscriptionRepo,
		subscriptionGraphqlRepo: repos.SubscriptionGraphqlRepo,
		usageChargeGraphqlRepo:  repos.UsageChargeGraphqlRepo,
	}
}

// CreateUsageSubscription 创建用量订阅
func (s *SubscriptionService) CreateUsageSubscription(
	ctx context.Context,
	userID int64,
	shopDomain string,
	planName string,
	cappedAmount decimal.Decimal,
	currency string,
	terms string,
	returnURL string,
	isTest bool,
) (*billingEntity.UserSubscription, string, error) {

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
		ReturnURL: returnURL,
	}

	// 2. 创建 Shopify 订阅
	subscription, confirmationURL, err := s.subscriptionGraphqlRepo.CreateSubscription(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Shopify subscription: %v", err)
	}

	// 3. 保存到本地数据库
	userSubscription := &billingEntity.UserSubscription{
		UserID:                 userID,
		ShopDomain:             shopDomain,
		SubscriptionID:         subscription.ID,
		SubscriptionName:       subscription.Name,
		SubscriptionStatus:     subscription.Status,
		SubscriptionLineItemID: subscription.LineItems[0].ID,
		PricingType:            billingEntity.PricingTypeRecurring,
		CappedAmount:           cappedAmount,
		Currency:               currency,
		BalanceUsed:            decimal.Zero,
		Terms:                  terms,
		CurrentPeriodStart:     time.Now().Unix(),
		CurrentPeriodEnd:       parseShopifyTime(subscription.CurrentPeriodEnd),
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
) (*billingEntity.UserSubscription, string, error) {

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
	subscription, confirmationURL, err := s.subscriptionGraphqlRepo.CreateSubscription(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Shopify subscription: %v", err)
	}

	// 3. 保存到本地数据库
	userSubscription := &billingEntity.UserSubscription{
		UserID:                 userID,
		ShopDomain:             shopDomain,
		SubscriptionID:         subscription.ID,
		SubscriptionName:       subscription.Name,
		SubscriptionStatus:     subscription.Status,
		SubscriptionLineItemID: subscription.LineItems[0].ID,
		PricingType:            billingEntity.PricingTypeRecurring,
		CappedAmount:           price, // 对于循环订阅，这里存储价格
		Currency:               currency,
		BalanceUsed:            decimal.Zero,
		Terms:                  fmt.Sprintf("Recurring charge - %s", interval),
		CurrentPeriodStart:     time.Now().Unix(),
		CurrentPeriodEnd:       parseShopifyTime(subscription.CurrentPeriodEnd),
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
func (s *SubscriptionService) SyncSubscriptionStatus(ctx context.Context, userID int64) error {
	// 1. 从 Shopify 获取当前订阅
	currentSubscription, err := s.subscriptionGraphqlRepo.GetCurrentSubscription(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current subscription from Shopify: %v", err)
	}

	// 2. 更新本地数据库
	userSubscription, err := s.userSubscriptionRepo.GetActiveSubscription(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user subscription from database: %v", err)
	}

	if userSubscription != nil {
		userSubscription.SubscriptionStatus = currentSubscription.Status
		userSubscription.CurrentPeriodEnd = parseShopifyTime(currentSubscription.CurrentPeriodEnd)
		userSubscription.LastSyncTime = time.Now().Unix()
		userSubscription.UpdateTime = time.Now().Unix()

		err = s.userSubscriptionRepo.UpsertUserSubscription(ctx, userSubscription)
		if err != nil {
			return fmt.Errorf("failed to update user subscription: %v", err)
		}
	}

	return nil
}

// parseShopifyTime 解析 Shopify 时间格式
func parseShopifyTime(timeStr string) int64 {
	if timeStr == "" {
		return 0
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0
	}

	return t.Unix()
}
