package billings

const (
	BillingPeriodSummaryTableName = "billing_period_summary"
)

type BillingPeriodSummary struct {
	Id                    int64   `xorm:"bigint UNSIGNED 'id' comment('ID') pk autoincr notnull " json:"id"`                                                // ID
	UserId                int64   `xorm:"bigint UNSIGNED 'user_id' comment('用户ID') notnull " json:"user_id"`                                                // 用户ID
	ShopDomain            string  `xorm:"varchar(100) 'shop_domain' comment('店铺域名') notnull " json:"shop_domain"`                                           // 店铺域名
	SubscriptionId        int64   `xorm:"not null default 0'subscription_id' comment('关联的订阅ID') notnull " json:"subscription_id"`                           // 关联的订阅ID
	BillingPeriodStart    int64   `xorm:"bigint UNSIGNED 'billing_period_start' comment('账单周期开始时间') notnull default 0 " json:"billing_period_start"`        // 账单周期开始时间
	BillingPeriodEnd      int64   `xorm:"bigint UNSIGNED 'billing_period_end' comment('账单周期结束时间') notnull default 0 " json:"billing_period_end"`            // 账单周期结束时间
	BillCycle             string  `xorm:"varchar(20) 'bill_cycle' comment('账单周期标识（YYYY-MM-DD）') notnull " json:"bill_cycle"`                                // 账单周期标识（YYYY-MM-DD）
	TotalCommissionAmount float64 `xorm:"decimal(12, 2) 'total_commission_amount' comment('周期总抽成金额') notnull default 0.00 " json:"total_commission_amount"` // 周期总抽成金额
	PendingAmount         float64 `xorm:"decimal(12, 2) 'pending_amount' comment('待付金额') notnull default 0.00 " json:"pending_amount"`                      // 待付金额
	PaidAmount            float64 `xorm:"decimal(12, 2) 'paid_amount' comment('已付金额') notnull default 0.00 " json:"paid_amount"`                            // 已付金额
	ErrorAmount           float64 `xorm:"decimal(12, 2) 'error_amount' comment('失败金额') notnull default 0.00 " json:"error_amount"`                          // 失败金额
	BillCount             int32   `xorm:"int 'bill_count' comment('账单数量') notnull default 0 " json:"bill_count"`                                            // 账单数量
	OrderCount            int32   `xorm:"int 'order_count' comment('订单数量') notnull default 0 " json:"order_count"`                                          // 订单数量
	Currency              string  `xorm:"varchar(10) 'currency' comment('货币类型') notnull " json:"currency"`                                                  // 货币类型
	ProtectifyType        string  `xorm:"varchar(30) 'protectify_type' comment('保险类型') notnull default general " json:"protectify_type"`                    // 保险类型
	SummaryStatus         string  `xorm:"varchar(20) 'summary_status' comment('周期状态：open-开放，closed-已关闭') notnull default open " json:"summary_status"`      // 周期状态：open-开放，closed-已关闭
	Remarks               string  `xorm:"varchar(255) 'remarks' comment('备注信息') notnull " json:"remarks"`                                                   // 备注信息
	TotalProtectifyAmount float64 `xorm:"decimal(12, 2) 'total_protectify_amount' comment('总保险金额') notnull default 0.00 " json:"total_protectify_amount"`   // 总保险金额
	TotalOrderAmount      float64 `xorm:"decimal(12, 2) 'total_order_amount' comment('总订单金额') notnull default 0.00 " json:"total_order_amount"`             // 总订单金额
	TotalRefundAmount     float64 `xorm:"decimal(12, 2) 'total_refund_amount' comment('总退款金额') notnull default 0.00 " json:"total_refund_amount"`           // 总退款金额
	BusinessMonth         string  `xorm:"varchar(7) 'business_month' comment('业务月份（YYYY-MM）') notnull " json:"business_month"`                              // 业务月份（YYYY-MM）
	IsTestPeriod          int8    `xorm:"tinyint 'is_test_period' comment('是否测试周期：0-否，1-是') notnull default 0 " json:"is_test_period"`                      // 是否测试周期：0-否，1-是
	Version               int32   `xorm:"int 'version' comment('版本号（用于乐观锁）') notnull default 1 " json:"version"`                                            // 版本号（用于乐观锁）
	LastSyncTime          int64   `xorm:"bigint UNSIGNED 'last_sync_time' comment('最后同步时间') notnull default 0 " json:"last_sync_time"`                      // 最后同步时间
	CreateTime            int64   `xorm:"created bigint UNSIGNED 'create_time' comment('创建时间') notnull " json:"create_time"`                                // 创建时间
	UpdateTime            int64   `xorm:"updated bigint UNSIGNED 'update_time' comment('修改时间') notnull " json:"update_time"`                                // 修改时间
}

func (s BillingPeriodSummary) TableName() string {
	return BillingPeriodSummaryTableName
}
