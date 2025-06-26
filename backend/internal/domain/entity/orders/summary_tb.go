package orders

import (
	"time"
)

// OrderSummary 用户订单记录统计
type OrderSummary struct {
	Id         int     `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid        int     `xorm:"'uid' int(11) notnull comment('用户ID')" json:"uid"`
	Today      int64   `xorm:"'today' int(11) default 0 comment('当天0点时间戳')" json:"today"`
	Orders     int     `xorm:"'orders' float default 0 comment('订单数')" json:"orders"`
	Sales      float64 `xorm:"'sales' float default 0 comment('销售金额')" json:"sales"`
	Refund     float64 `xorm:"'refund' float default 0 comment('退款金额')" json:"refund"`
	CreateTime int64   `xorm:"'create_time' bigint(20) notnull comment('创建时间')" json:"create_time"`
	UpdateTime int64   `xorm:"'update_time' bigint(20) notnull comment('修改时间')" json:"update_time"`
}

//func (m *UserProduct) TableName() string {
//	return "in_user_product"
//}

func (o *OrderSummary) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	o.CreateTime = now
	o.UpdateTime = now
}

func (o *OrderSummary) BeforeUpdate() {
	now := time.Now().Unix()
	o.UpdateTime = now
}
