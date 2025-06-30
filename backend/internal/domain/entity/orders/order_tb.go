package orders

// UserOrder 用户订单主表
type UserOrder struct {
	Id                int64   `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	UserID            int64   `xorm:"'user_id' bigint(20) notnull comment('用户id')" json:"user_id"`
	OrderId           int64   `xorm:"'order_id' bigint(20) notnull default '' comment('Shopify订单ID')" json:"order_id"`
	OrderName         string  `xorm:"'order_name' varchar(50) notnull default '' comment('订单编号（#xxx）')" json:"order_name"`
	OrderCreatedAt    int64   `xorm:"'order_created_at' bigint(20) notnull default 0 comment('订单创建时间')" json:"order_created_at"`
	OrderCompletionAt int64   `xorm:"'order_completion_at' bigint(20) notnull default 0 comment('订单完成时间')" json:"order_completion_at"`
	FinancialStatus   string  `xorm:"'financial_status' varchar(50) notnull default '' comment('支付状态')" json:"financial_status"`
	TotalPriceAmount  float64 `xorm:"'total_price_amount' decimal(12,2) notnull default 0.00 comment('订单总金额')" json:"total_price_amount"`
	RefundPriceAmount float64 `xorm:"'refund_price_amount' decimal(12,2) notnull default 0.00 comment('退款总金额')" json:"refund_price_amount"`
	InsuranceAmount   float64 `xorm:"'insurance_amount' decimal(12,2) notnull default 0.00 comment('保险金额')" json:"insurance_amount"`
	Currency          string  `xorm:"'currency' varchar(10) notnull default '' comment('货币类型')" json:"currency"`
	SkuNum            int     `xorm:"'sku_num' int(11) notnull default 0 comment('sku购买数量')" json:"sku_num"`
	IsDel             int     `xorm:"'is_del' tinyint(1) notnull default 0 comment('删除状态 0 正常 1 已删除')" json:"is_del"`
	CreateTime        int64   `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime        int64   `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
