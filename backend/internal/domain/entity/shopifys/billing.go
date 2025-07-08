package shopifys

import "github.com/shopspring/decimal"

// AppUsageRecordCreateResponse appUsageRecordCreate 响应结构
type AppUsageRecordCreateResponse struct {
	AppUsageRecord *AppUsageRecord `json:"appUsageRecord"`
	UserErrors     []UserError     `json:"userErrors"`
}

// AppUsageRecord 用量记录
type AppUsageRecord struct {
	ID                   string                   `json:"id"`
	Description          string                   `json:"description"`
	Price                MoneyV2                  `json:"price"`
	CreatedAt            string                   `json:"createdAt"`
	SubscriptionLineItem *AppSubscriptionLineItem `json:"subscriptionLineItem,omitempty"`
}

// MoneyV2 货币类型
type MoneyV2 struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

// AppSubscriptionLineItem 订阅项目
type AppSubscriptionLineItem struct {
	ID   string `json:"id"`
	Plan struct {
		PricingDetails AppUsagePricingDetails `json:"pricingDetails"`
	} `json:"plan"`
}

// AppUsagePricingDetails 用量定价详情
type AppUsagePricingDetails struct {
	CappedAmount MoneyV2 `json:"cappedAmount"`
	Terms        string  `json:"terms"`
	BalanceUsed  MoneyV2 `json:"balanceUsed"`
}

// UserError 用户错误
type UserError struct {
	Field   []string `json:"field"`
	Message string   `json:"message"`
}

// AppSubscriptionCreateInput 创建订阅输入
type AppSubscriptionCreateInput struct {
	Name      string                         `json:"name"`
	Test      bool                           `json:"test"`
	TrialDays int                            `json:"trialDays,omitempty"`
	LineItems []AppSubscriptionLineItemInput `json:"lineItems"`
	ReturnURL string                         `json:"returnUrl"`
}
type AppSubscriptionLineItemInput struct {
	Plan AppPlanInput `json:"plan"`
}

// AppPlanInput 计划输入
type AppPlanInput struct {
	AppUsagePricingDetails     *AppUsagePricingDetailsInput     `json:"appUsagePricingDetails,omitempty"`
	AppRecurringPricingDetails *AppRecurringPricingDetailsInput `json:"appRecurringPricingDetails,omitempty"`
}

// MoneyInput 金额设置
type MoneyInput struct {
	Amount       decimal.Decimal `json:"amount"`
	CurrencyCode string          `json:"currencyCode"`
}

// AppUsagePricingDetailsInput 用量定价详情输入
type AppUsagePricingDetailsInput struct {
	CappedAmount MoneyInput `json:"cappedAmount"`
	Terms        string     `json:"terms"`
}

// AppRecurringPricingDetailsInput 循环定价详情输入
type AppRecurringPricingDetailsInput struct {
	Price    MoneyInput `json:"price"`
	Interval string     `json:"interval"` // EVERY_30_DAYS, ANNUAL
}

// AppSubscriptionCreateResponse 创建订阅响应
type AppSubscriptionCreateResponse struct {
	AppSubscription *AppSubscription `json:"appSubscription"`
	ConfirmationURL string           `json:"confirmationUrl"`
	UserErrors      []UserError      `json:"userErrors"`
}

// AppSubscription 订阅信息
type AppSubscription struct {
	ID               string                            `json:"id"`
	Name             string                            `json:"name"`
	Status           string                            `json:"status"`
	Test             bool                              `json:"test"`
	CreatedAt        string                            `json:"createdAt"`
	CurrentPeriodEnd string                            `json:"currentPeriodEnd"`
	TrialDays        int                               `json:"trialDays"`
	LineItems        []AppSubscriptionLineItemResponse `json:"lineItems"`
}

// AppSubscriptionLineItemResponse 订阅项目响应
type AppSubscriptionLineItemResponse struct {
	ID   string `json:"id"`
	Plan struct {
		PricingDetails interface{} `json:"pricingDetails"`
	} `json:"plan"`
}
