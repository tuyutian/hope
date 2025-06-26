package orders

import (
	"time"
)

// UserOrder 用户订单主表
type UserOrder struct {
	Id                int     `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid               int     `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	OrderId           string  `xorm:"'order_id' varchar(50) comment('Shopify订单ID')" json:"order_id"`
	OrderName         string  `xorm:"'order_name' varchar(50) comment('订单编号（#xxx）')" json:"order_name"`
	OrderCreatedAt    int64   `xorm:"'order_created_at' int(11) default 0 comment('订单创建时间')" json:"order_created_at"`
	OrderCompletionAt int64   `xorm:"'order_completion_at' int(11) default 0 comment('订单完成时间')" json:"order_completion_at"`
	FinancialStatus   string  `xorm:"'financial_status' varchar(50) comment('支付状态')" json:"financial_status"`
	TotalPriceAmount  float64 `xorm:"'total_price_amount' decimal(10,2) default 0.00 comment('订单总金额')" json:"total_price_amount"`
	RefundPriceAmount float64 `xorm:"'refund_price_amount' decimal(10,2) default 0.00 comment('退款总金额')" json:"refund_price_amount"`
	InsuranceAmount   float64 `xorm:"'insurance_amount' decimal(10,2) default 0.00 comment('保险金额')" json:"insurance_amount"`
	Currency          string  `xorm:"'currency' varchar(10) comment('货币类型')" json:"currency"`
	SkuNum            int     `xorm:"'sku_num' int(11) default 0 comment('sku购买数量')" json:"sku_num"`
	IsDel             int     `xorm:"'is_del' tinyint(1) default 0 notnull comment('删除状态 0 正常 1 已删除')" json:"is_del"`
	CreateTime        int64   `xorm:"'create_time' int(11) default 0 comment('创建时间')" json:"create_time"`
	UpdateTime        int64   `xorm:"'update_time' int(11) default 0 comment('修改时间')" json:"update_time"`
}

// OrderStatistics 用于存储查询结果
type OrderStatistics struct {
	TotalRefund    float64 `xorm:"'total_refund'"`
	TotalInsurance float64 `xorm:"'total_insurance'"`
	TotalOrders    int     `xorm:"'total_orders'"`
}

//func (j *JobProduct) TableName() string {
//	return "in_job_product"
//}

func (o *UserOrder) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	o.CreateTime = now
	o.UpdateTime = now
}

func (o *UserOrder) BeforeUpdate() {
	now := time.Now().Unix()
	o.UpdateTime = now
}
