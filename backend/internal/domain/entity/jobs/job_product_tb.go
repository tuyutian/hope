package jobs

import (
	"time"
)

// JobProduct UserProduct 用户上传记录表
type JobProduct struct {
	Id            int   `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid           int   `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	UserProductId int   `xorm:"'user_product_id' int(11) notnull comment('用户产品ID')" json:"user_product_id"`
	JobTime       int64 `xorm:"'job_time' int(11) notnull comment('队列时间')" json:"job_time"`
	IsSuccess     int   `xorm:"'is_success' tinyint(1) default 0 notnull comment('处理状态 0 未处理完成 1 处理成功')" json:"is_success"`
	CreateTime    int64 `xorm:"'create_time' int(11) notnull comment('创建时间')" json:"create_time"`
	UpdateTime    int64 `xorm:"'update_time' int(11) notnull comment('修改时间')" json:"update_time"`
}

//func (j *JobProduct) TableName() string {
//	return "in_job_product"
//}

func (j *JobProduct) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	j.CreateTime = now
	j.UpdateTime = now
}

func (j *JobProduct) BeforeUpdate() {
	now := time.Now().Unix()
	j.UpdateTime = now
}
