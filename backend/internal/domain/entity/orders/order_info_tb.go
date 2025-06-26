package orders

import (
	"time"
)

// UserOrderInfo 订单详情表
type UserOrderInfo struct {
	Id              int     `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid             int     `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	UserOrderId     int     `xorm:"'user_order_id' int(11) notnull comment('主表ID')" json:"user_order_id"`
	Sku             string  `xorm:"'sku' varchar(100) comment('SKU')" json:"sku"`
	VariantId       string  `xorm:"'variant_id' varchar(100) comment('变体ID')" json:"variant_id"`
	VariantTitle    string  `xorm:"'variant_title' varchar(255) comment('变体标题')" json:"product_title"`
	Quantity        int     `xorm:"'quantity' int(11) default 0 comment('购买数量')" json:"quantity"`
	UnitPriceAmount float64 `xorm:"'unit_price_amount' decimal(10,2) default 0.00 comment('单价金额')" json:"unit_price_amount"`
	Currency        string  `xorm:"'currency' varchar(10) comment('货币类型')" json:"currency"`
	RefundNum       int     `xorm:"'refund_num' int(11) default 0 comment('退款数量')" json:"refund_num"`
	IsInsurance     int     `xorm:"'is_insurance' tinyint(1) default 0 comment('是否是保险产品')" json:"is_insurance"`
	CreateTime      int64   `xorm:"'create_time' int(11) default 0 comment('创建时间')" json:"create_time"`
	UpdateTime      int64   `xorm:"'update_time' int(11) default 0 comment('修改时间')" json:"update_time"`
}

//func (j *JobProduct) TableName() string {
//	return "in_job_product"
//}

func (o *UserOrderInfo) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	o.CreateTime = now
	o.UpdateTime = now
}

func (o *UserOrderInfo) BeforeUpdate() {
	now := time.Now().Unix()
	o.UpdateTime = now
}
