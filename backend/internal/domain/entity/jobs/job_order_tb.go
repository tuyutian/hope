package jobs

// JobOrder 用户订单同步记录表
type JobOrder struct {
	Id         int64 `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	OrderId    int64 `xorm:"'order_id' bigint(20) notnull default 0 comment('shopify 订单id')" json:"order_id"`
	UserID     int64 `xorm:"'user_id' bigint(20) notnull default 0 comment('店铺')" json:"user_id"`
	JobTime    int64 `xorm:"'job_time' bigint(20) notnull comment('队列时间(毫秒时间戳)')" json:"job_time"`
	IsSuccess  int   `xorm:"'is_success' tinyint(1) default 0 notnull comment('处理状态 0 未处理完成 1 处理成功')" json:"is_success"`
	CreateTime int64 `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime int64 `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
