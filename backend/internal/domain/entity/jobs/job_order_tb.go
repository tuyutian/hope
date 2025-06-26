package jobs

import (
	"time"
)

// JobOrder 用户订单同步记录表
type JobOrder struct {
	Id         int    `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	OrderId    string `xorm:"'order_id' varchar(100) comment('shopify 订单id')" json:"order_id"`
	Name       string `xorm:"'name' varchar(100) comment('店铺name')" json:"name"`
	JobTime    int    `xorm:"'job_time' int(11) notnull comment('队列时间')" json:"job_time"`
	IsSuccess  int    `xorm:"'is_success' tinyint(1) default 0 notnull comment('处理状态 0 未处理完成 1 处理成功')" json:"is_success"`
	CreateTime int64  `xorm:"'create_time' int(11) notnull comment('创建时间')" json:"create_time"`
	UpdateTime int64  `xorm:"updated 'update_time' int(11) notnull comment('修改时间')" json:"update_time"`
}

//func (j *JobProduct) TableName() string {
//	return "in_job_product"
//}

func (j *JobOrder) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	j.CreateTime = now
	j.UpdateTime = now
}

func (j *JobOrder) BeforeUpdate() {
	now := time.Now().Unix()
	j.UpdateTime = now
}
