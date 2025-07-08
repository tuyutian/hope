package billings

import "github.com/shopspring/decimal"

// 扣费状态
const (
	ChargeStatusPending = 0 // 待提交
	ChargeStatusCharged = 1 // 已提交到Shopify
	ChargeStatusFailed  = 2 // 提交失败
)

type CommissionBill struct {
	Id                   int64           `json:"id" xorm:"pk autoincr 'id' comment('ID')"`
	ChargeId             int64           `json:"charge_id" xorm:"notnull 'charge_id' comment('账单编号')"`
	UserId               int64           `json:"user_id" xorm:"notnull 'user_id' comment('用户ID')"`
	UserOrderId          int64           `json:"user_order_id" xorm:"notnull 'user_order_id' comment('关联的订单ID')"`
	OrderName            string          `json:"order_name" xorm:"varchar(50) notnull default '' 'order_name' comment('Shopify订单编号')"`
	BillCycle            string          `json:"bill_cycle" xorm:"varchar(20) notnull 'bill_cycle' comment('账单周期（YYYY-MM）')"`
	CommissionAmount     decimal.Decimal `json:"commission_amount" xorm:"decimal(12,2) notnull default 0.00 'commission_amount' comment('抽成金额')"`
	Currency             string          `json:"currency" xorm:"varchar(10) notnull default '' 'currency' comment('货币类型')"`
	ShopifyUsageRecordId string          `json:"shopify_usage_record_id" xorm:"varchar(100) notnull default '' 'shopify_usage_record_id' comment('Shopify用量记录ID')"`
	ChargeStatus         int             `json:"charge_status" xorm:"tinyint notnull default 0 'charge_status' comment('扣费状态：0-待提交, 1-已提交, 2-提交失败')"`
	ErrorMessage         string          `json:"error_message" xorm:"text 'error_message' comment('错误信息')"`
	ChargedAt            int64           `json:"charged_at" xorm:"notnull default 0 'charged_at' comment('扣费时间')"`
	CreateTime           int64           `json:"create_time" xorm:"updated notnull 'create_time' comment('创建时间')"`
	UpdateTime           int64           `json:"update_time" xorm:"updated notnull 'update_time' comment('修改时间')"`
}

func (c *CommissionBill) TableName() string {
	return "commission_bill"
}
