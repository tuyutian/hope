package shopifys

import "time"

type Token struct {
	Token string `json:"access_token"`
	Scope string `json:"scope"`
}

// SubscriptionWebhookPayload 订阅 webhook 载荷
type SubscriptionWebhookPayload struct {
	AdminGraphqlApiId     string    `json:"admin_graphql_api_id"`
	Name                  string    `json:"name"`
	Status                string    `json:"status"`
	AdminGraphqlApiShopId string    `json:"admin_graphql_api_shop_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Currency              string    `json:"currency"`
	CappedAmount          string    `json:"capped_amount"`
	Price                 string    `json:"price"`
	Interval              string    `json:"interval"`
	PlanHandle            string    `json:"plan_handle"`
}
type SubscriptionApproachingPayload struct {
	AdminGraphqlApiId     string    `json:"admin_graphql_api_id"`
	Name                  string    `json:"name"`
	BalanceUsed           int       `json:"balance_used"`
	CappedAmount          string    `json:"capped_amount"`
	CurrencyCode          string    `json:"currency_code"`
	AdminGraphqlApiShopId string    `json:"admin_graphql_api_shop_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// LineItem 订阅项目
type LineItem struct {
	ID   string      `json:"id"`
	Plan interface{} `json:"plan"`
}

// UsageRecordWebhookPayload 使用记录 webhook 载荷
type UsageRecordWebhookPayload struct {
	ID                string    `json:"id"`
	SubscriptionAppId string    `json:"subscription_app_id"`
	Description       string    `json:"description"`
	Price             string    `json:"price"`
	CreatedAt         time.Time `json:"created_at"`
	AdminGraphqlApiId string    `json:"admin_graphql_api_id"`
}

// PaymentCallbackRequest 通用支付回调请求
type PaymentCallbackRequest struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Amount      string                 `json:"amount"`
	Currency    string                 `json:"currency"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
