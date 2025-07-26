package shopifys

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

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
	ID               string                    `json:"id"`
	Name             string                    `json:"name"`
	Status           string                    `json:"status"`
	Test             bool                      `json:"test"`
	CreatedAt        string                    `json:"createdAt"`
	CurrentPeriodEnd string                    `json:"currentPeriodEnd"`
	ReturnUrl        string                    `json:"returnUrl"`
	TrialDays        int                       `json:"trialDays"`
	LineItems        []AppSubscriptionLineItem `json:"lineItems"`
}

// AppSubscriptionLineItem 订阅项目响应
type AppSubscriptionLineItem struct {
	ID   string              `json:"id"`
	Plan AppSubscriptionPlan `json:"plan"`
}

// AppSubscriptionPlan 订阅计划响应
type AppSubscriptionPlan struct {
	PricingDetails AppSubscriptionPricingDetails `json:"pricingDetails"`
}

// UnmarshalJSON 自定义 UnmarshalJSON 方法来处理 Union 类型
func (p *AppSubscriptionPlan) UnmarshalJSON(data []byte) error {
	// 先解析到一个临时结构体
	var temp struct {
		PricingDetails json.RawMessage `json:"pricingDetails"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// 尝试解析为不同的定价类型
	// 首先检查是否包含使用量定价的字段
	var usagePricing AppUsagePricing
	if err := json.Unmarshal(temp.PricingDetails, &usagePricing); err == nil {
		// 检查是否包含使用量定价的特有字段
		var tempCheck map[string]interface{}
		json.Unmarshal(temp.PricingDetails, &tempCheck)
		if _, hasCappedAmount := tempCheck["cappedAmount"]; hasCappedAmount {
			if _, hasTerms := tempCheck["terms"]; hasTerms {
				p.PricingDetails = usagePricing
				return nil
			}
		}
	}

	// 尝试解析为循环定价
	var recurringPricing AppRecurringPricing
	if err := json.Unmarshal(temp.PricingDetails, &recurringPricing); err == nil {
		// 检查是否包含循环定价的特有字段
		var tempCheck map[string]interface{}
		json.Unmarshal(temp.PricingDetails, &tempCheck)
		if _, hasPrice := tempCheck["price"]; hasPrice {
			if _, hasInterval := tempCheck["interval"]; hasInterval {
				p.PricingDetails = recurringPricing
				return nil
			}
		}
	}

	return fmt.Errorf("unknown pricing details type")
}

// AppSubscriptionPricingDetails 订阅定价详情（Union 类型）
type AppSubscriptionPricingDetails interface {
	GetPricingType() string
}

// AppRecurringPricing 循环计费定价
type AppRecurringPricing struct {
	Discount   *AppSubscriptionDiscount `json:"discount,omitempty"`
	Interval   string                   `json:"interval"`
	Price      MoneyV2                  `json:"price"`
	PlanHandle string                   `json:"planHandle"`
}

func (a AppRecurringPricing) GetPricingType() string {
	return "AppRecurringPricing"
}

// AppUsagePricing 使用量计费定价
type AppUsagePricing struct {
	BalanceUsed  MoneyV2 `json:"balanceUsed"`
	CappedAmount MoneyV2 `json:"cappedAmount"`
	Interval     string  `json:"interval"`
	Terms        string  `json:"terms"`
}

func (a AppUsagePricing) GetPricingType() string {
	return "AppUsagePricing"
}

// AppSubscriptionDiscount 订阅折扣
type AppSubscriptionDiscount struct {
	DurationLimitInIntervals     *int                         `json:"durationLimitInIntervals,omitempty"`
	RemainingDurationInIntervals *int                         `json:"remainingDurationInIntervals,omitempty"`
	Value                        AppSubscriptionDiscountValue `json:"value"`
}

// UnmarshalJSON 为 AppSubscriptionDiscount 添加自定义反序列化
func (d *AppSubscriptionDiscount) UnmarshalJSON(data []byte) error {
	// 先解析基本字段
	var temp struct {
		DurationLimitInIntervals     *int            `json:"durationLimitInIntervals,omitempty"`
		RemainingDurationInIntervals *int            `json:"remainingDurationInIntervals,omitempty"`
		Value                        json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	d.DurationLimitInIntervals = temp.DurationLimitInIntervals
	d.RemainingDurationInIntervals = temp.RemainingDurationInIntervals

	// 尝试解析 Value 字段
	// 首先尝试百分比折扣
	var percentageDiscount AppSubscriptionDiscountPercentage
	if err := json.Unmarshal(temp.Value, &percentageDiscount); err == nil {
		var tempCheck map[string]interface{}
		json.Unmarshal(temp.Value, &tempCheck)
		if _, hasPercentage := tempCheck["percentage"]; hasPercentage {
			d.Value = percentageDiscount
			return nil
		}
	}

	// 尝试金额折扣
	var amountDiscount AppSubscriptionDiscountAmount
	if err := json.Unmarshal(temp.Value, &amountDiscount); err == nil {
		var tempCheck map[string]interface{}
		json.Unmarshal(temp.Value, &tempCheck)
		if _, hasAmount := tempCheck["amount"]; hasAmount {
			d.Value = amountDiscount
			return nil
		}
	}

	return fmt.Errorf("unknown discount value type")
}

// AppSubscriptionDiscountValue 折扣值（Union 类型）
type AppSubscriptionDiscountValue interface {
	GetDiscountType() string
}

// AppSubscriptionDiscountAmount 固定金额折扣
type AppSubscriptionDiscountAmount struct {
	Amount MoneyV2 `json:"amount"`
}

func (a AppSubscriptionDiscountAmount) GetDiscountType() string {
	return "AppSubscriptionDiscountAmount"
}

// AppSubscriptionDiscountPercentage 百分比折扣
type AppSubscriptionDiscountPercentage struct {
	Percentage float64 `json:"percentage"`
}

func (a AppSubscriptionDiscountPercentage) GetDiscountType() string {
	return "AppSubscriptionDiscountPercentage"
}

// MoneyV2 货币金额
type MoneyV2 struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

// GetRecurringPricing 获取循环计费信息
func (a *AppSubscription) GetRecurringPricing() (*AppRecurringPricing, error) {
	if len(a.LineItems) == 0 {
		return nil, fmt.Errorf("no line items found")
	}

	pricing := a.LineItems[0].Plan.PricingDetails
	if recurring, ok := pricing.(AppRecurringPricing); ok {
		return &recurring, nil
	}

	return nil, fmt.Errorf("not a recurring pricing subscription")
}

// GetUsagePricing 获取使用量计费信息
func (a *AppSubscription) GetUsagePricing() (*AppUsagePricing, error) {
	if len(a.LineItems) == 0 {
		return nil, fmt.Errorf("no line items found")
	}

	pricing := a.LineItems[0].Plan.PricingDetails
	if usage, ok := pricing.(AppUsagePricing); ok {
		return &usage, nil
	}

	return nil, fmt.Errorf("not a usage pricing subscription")
}

// IsRecurringSubscription 判断是否为循环订阅
func (a *AppSubscription) IsRecurringSubscription() bool {
	_, err := a.GetRecurringPricing()
	return err == nil
}

// IsUsageSubscription 判断是否为使用量订阅
func (a *AppSubscription) IsUsageSubscription() bool {
	_, err := a.GetUsagePricing()
	return err == nil
}

// GetSubscriptionPrice 获取订阅价格
func (a *AppSubscription) GetSubscriptionPrice() (*MoneyV2, error) {
	if recurring, err := a.GetRecurringPricing(); err == nil {
		return &recurring.Price, nil
	}

	if usage, err := a.GetUsagePricing(); err == nil {
		return &usage.CappedAmount, nil
	}

	return nil, fmt.Errorf("unable to determine subscription price")
}

// GetDiscountAmount 获取折扣金额
func (d *AppSubscriptionDiscount) GetDiscountAmount() (*MoneyV2, error) {
	if amount, ok := d.Value.(AppSubscriptionDiscountAmount); ok {
		return &amount.Amount, nil
	}
	return nil, fmt.Errorf("discount is not amount-based")
}

// GetDiscountPercentage 获取折扣百分比
func (d *AppSubscriptionDiscount) GetDiscountPercentage() (float64, error) {
	if percentage, ok := d.Value.(AppSubscriptionDiscountPercentage); ok {
		return percentage.Percentage, nil
	}
	return 0, fmt.Errorf("discount is not percentage-based")
}

// HasDiscount 判断是否有折扣
func (a AppRecurringPricing) HasDiscount() bool {
	return a.Discount != nil
}
