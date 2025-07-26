package orders

// UserOrderInfo 订单详情表
type UserOrderInfo struct {
	Id              int64   `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	UserID          int64   `xorm:"'user_id' bigint(20) notnull comment('用户id')" json:"user_id"`
	UserOrderId     int64   `xorm:"'user_order_id' bigint(20) notnull comment('主表ID')" json:"user_order_id"`
	Sku             string  `xorm:"'sku' varchar(100) notnull default '' comment('SKU')" json:"sku"`
	VariantId       int64   `xorm:"'variant_id' bigint(20) notnull default '' comment('变体ID')" json:"variant_id"`
	VariantTitle    string  `xorm:"'variant_title' varchar(255) notnull default '' comment('变体标题')" json:"variant_title"`
	Quantity        int     `xorm:"'quantity' int(11) notnull default 0 comment('购买数量')" json:"quantity"`
	UnitPriceAmount float64 `xorm:"'unit_price_amount' decimal(12,2) notnull default 0.00 comment('单价金额')" json:"unit_price_amount"`
	Currency        string  `xorm:"'currency' varchar(10) notnull default '' comment('货币类型')" json:"currency"`
	RefundNum       int     `xorm:"'refund_num' int(11) notnull default 0 comment('退款数量')" json:"refund_num"`
	IsInsurance     int     `xorm:"'is_insurance' tinyint(1) notnull default 0 comment('是否是保险产品')" json:"is_insurance"`
	CreateTime      int64   `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime      int64   `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
