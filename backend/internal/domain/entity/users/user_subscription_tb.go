package users

import (
	"github.com/shopspring/decimal"
)

// UserSubscription 用户订阅信息表
type UserSubscription struct {
	ID                     int64           `xorm:"pk autoincr 'id'" json:"id"`
	UserID                 int64           `xorm:"not null 'user_id'" json:"user_id"`
	ShopDomain             string          `xorm:"varchar(100) not null default '' 'shop_domain'" json:"shop_domain"`
	SubscriptionID         int64           `xorm:"not null default 0 'subscription_id'" json:"subscription_id"`
	SubscriptionName       string          `xorm:"varchar(100) not null default '' 'subscription_name'" json:"subscription_name"`
	SubscriptionStatus     string          `xorm:"varchar(20) not null default '' 'subscription_status'" json:"subscription_status"`
	SubscriptionLineItemID int64           `xorm:"not null default 0 'subscription_line_item_id'" json:"subscription_line_item_id"`
	PricingType            string          `xorm:"varchar(20) not null default '' 'pricing_type'" json:"pricing_type"`
	CappedAmount           decimal.Decimal `xorm:"decimal(12,2) not null default 0.00 'capped_amount'" json:"capped_amount"`
	Currency               string          `xorm:"varchar(10) not null default '' 'currency'" json:"currency"`
	BalanceUsed            decimal.Decimal `xorm:"decimal(12,2) not null default 0.00 'balance_used'" json:"balance_used"`
	Terms                  string          `xorm:"text 'terms'" json:"terms"`
	CurrentPeriodStart     int64           `xorm:"not null default 0 'current_period_start'" json:"current_period_start"`
	CurrentPeriodEnd       int64           `xorm:"not null default 0 'current_period_end'" json:"current_period_end"`
	TrialDays              int             `xorm:"not null default 0 'trial_days'" json:"trial_days"`
	TestSubscription       bool            `xorm:"tinyint not null default 0 'test_subscription'" json:"test_subscription"`
	AppInstallationID      string          `xorm:"varchar(100) not null default '' 'app_installation_id'" json:"app_installation_id"`
	LastSyncTime           int64           `xorm:"not null default 0 'last_sync_time'" json:"last_sync_time"`
	CreateTime             int64           `xorm:"created not null 'create_time'" json:"create_time"`
	UpdateTime             int64           `xorm:"updated not null 'update_time'" json:"update_time"`
}

// TableName 指定表名
func (u *UserSubscription) TableName() string {
	return "user_subscription"
}

// 订阅状态常量
const (
	SubscriptionStatusActive    = "ACTIVE"
	SubscriptionStatusCancelled = "CANCELLED"
	SubscriptionStatusDeclined  = "DECLINED"
	SubscriptionStatusExpired   = "EXPIRED"
	SubscriptionStatusFrozen    = "FROZEN"
	SubscriptionStatusPending   = "PENDING"
)

// 定价类型常量
const (
	PricingTypeAnnual    = "ANNUAL"
	PricingTypeRecurring = "RECURRING"
	PricingTypeOneTime   = "ONE_TIME"
)
