package orders

// OrderSummary 用户订单记录统计
type OrderSummary struct {
	Id         int64   `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	UserID     int64   `xorm:"'user_id' bigint(20) notnull comment('用户ID')" json:"user_id"`
	Today      int64   `xorm:"'today' bigint(20) notnull default 0 comment('当天0点时间戳')" json:"today"`
	Orders     int     `xorm:"'orders' int(11) notnull default 0 comment('订单数')" json:"orders"`
	Sales      float64 `xorm:"'sales' decimal(12,2) notnull default 0.00 comment('销售金额')" json:"sales"`
	Refund     float64 `xorm:"'refund' decimal(12,2) notnull default 0.00 comment('退款金额')" json:"refund"`
	CreateTime int64   `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime int64   `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
