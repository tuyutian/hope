package billings

// 扣费状态
const (
	ChargeStatusPending = 0 // 待提交
	ChargeStatusCharged = 1 // 已提交到Shopify
	ChargeStatusFailed  = 2 // 提交失败
)

type CommissionBill struct {
	Id                    int64   `xorm:"bigint UNSIGNED 'id' comment('ID') pk autoincr notnull " json:"id"`                                                                     // ID
	ChargeId              int64   `xorm:"bigint UNSIGNED 'charge_id' comment('账单编号') notnull " json:"charge_id"`                                                                 // 账单编号
	UserId                int64   `xorm:"bigint UNSIGNED 'user_id' comment('用户ID') notnull " json:"user_id"`                                                                     // 用户ID
	UserOrderId           int64   `xorm:"bigint UNSIGNED 'user_order_id' comment('关联的订单ID') notnull " json:"user_order_id"`                                                      // 关联的订单ID
	OrderName             string  `xorm:"varchar(50) 'order_name' comment('Shopify订单编号') notnull " json:"order_name"`                                                            // Shopify订单编号
	BillingPeriodStart    int64   `xorm:"bigint UNSIGNED 'billing_period_start' comment('账单周期开始时间') notnull default 0 " json:"billing_period_start"`                             // 账单周期开始时间
	BillingPeriodEnd      int64   `xorm:"bigint UNSIGNED 'billing_period_end' comment('账单周期结束时间') notnull default 0 " json:"billing_period_end"`                                 // 账单周期结束时间
	BillCycle             string  `xorm:"varchar(20) 'bill_cycle' comment('账单周期标识（YYYY-MM-DD）') notnull " json:"bill_cycle"`                                                     // 账单周期标识（YYYY-MM-DD）
	CommissionAmount      float64 `xorm:"decimal(12, 2) 'commission_amount' comment('抽成金额') notnull default 0.00 " json:"commission_amount"`                                     // 抽成金额
	CommissionRate        float64 `xorm:"decimal(5, 2) 'commission_rate' comment('抽成比例（百分比）') notnull default 0.00 " json:"commission_rate"`                                     // 抽成比例（百分比）
	ProtectifyType        string  `xorm:"varchar(30) 'protectify_type' comment('保险类型：general-通用保险，product-产品保险，shipping-运输保险') notnull default general " json:"protectify_type"` // 保险类型：general-通用保险，product-产品保险，shipping-运输保险
	SubscriptionId        int64   `xorm:"bigint UNSIGNED 'subscription_id' comment('关联的订阅ID') notnull default 0 " json:"subscription_id"`                                        // 关联的订阅ID
	OrderProtectifyAmount float64 `xorm:"decimal(12, 2) 'order_protectify_amount' comment('订单保险金额') notnull default 0.00 " json:"order_protectify_amount"`                       // 订单保险金额
	OrderTotalAmount      float64 `xorm:"decimal(12, 2) 'order_total_amount' comment('订单总金额') notnull default 0.00 " json:"order_total_amount"`                                  // 订单总金额
	CommissionItems       string  `xorm:"text 'commission_items' comment('抽成明细项（JSON格式，包含保险项目等）') " json:"commission_items"`                                                     // 抽成明细项（JSON格式，包含保险项目等）
	Currency              string  `xorm:"varchar(10) 'currency' comment('货币类型') notnull " json:"currency"`                                                                       // 货币类型
	ShopifyUsageRecordId  string  `xorm:"varchar(100) 'shopify_usage_record_id' comment('Shopify用量记录ID') notnull " json:"shopify_usage_record_id"`                               // Shopify用量记录ID
	ChargeStatus          int8    `xorm:"tinyint 'charge_status' comment('扣费状态：0-待提交, 1-已提交, 2-提交失败') notnull default 0 " json:"charge_status"`                                  // 扣费状态：0-待提交, 1-已提交, 2-提交失败
	ErrorMessage          string  `xorm:"text 'error_message' comment('错误信息') " json:"error_message"`                                                                            // 错误信息
	ChargedAt             int64   `json:"charged_at" xorm:"notnull default 0 'charged_at' comment('扣费时间')"`
	CreateTime            int64   `json:"create_time" xorm:"created notnull 'create_time' comment('创建时间')"`
	UpdateTime            int64   `json:"update_time" xorm:"updated notnull 'update_time' comment('修改时间')"`
}

func (c CommissionBill) TableName() string {
	return "commission_bill"
}
