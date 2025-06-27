package jobs

// JobProduct UserProduct 用户上传记录表
type JobProduct struct {
	Id            int64 `xorm:"pk autoincr 'id' bigint(20) comment('ID')" json:"id"`
	UserID        int64 `xorm:"'user_id' bigint(20) notnull comment('用户id')" json:"user_id"`
	UserProductId int64 `xorm:"'user_product_id' bigint(20) notnull comment('用户产品ID')" json:"user_product_id"`
	JobTime       int64 `xorm:"'job_time' bigint(20) notnull comment('队列时间(毫秒时间戳)')" json:"job_time"`
	IsSuccess     int   `xorm:"'is_success' tinyint(1) default 0 notnull comment('处理状态 0 未处理完成 1 处理成功')" json:"is_success"`
	CreateTime    int64 `xorm:"created 'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime    int64 `xorm:"updated 'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}
