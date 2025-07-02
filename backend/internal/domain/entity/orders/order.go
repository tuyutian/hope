package orders

type QueryOrderEntity struct {
	UserID   int64  `json:"user_id"`
	Page     int    `form:"page" binding:"required,gte=1"` // gte=1 表示必须 ≥1
	PageSize int    `form:"page_size" binding:"required,gte=1,lte=50"`
	Type     int    `form:"type" binding:"required,oneof=0 1 2"`
	Query    string `form:"query" binding:"omitempty,max=20"`
}

// OrderStatistics 用于存储查询结果
type OrderStatistics struct {
	TotalRefund    float64 `xorm:"'total_refund'" json:"total_refund"`
	TotalInsurance float64 `xorm:"'total_insurance'" json:"total_insurance"`
	TotalOrders    int     `xorm:"'total_orders'" json:"total_orders"`
}

type OrderDashboardReq struct {
	UserID int64 `json:"user_id"`
	Days   int   `form:"days" binding:"required,oneof=30 90 365"`
}

type OrderWebHookReq struct {
	Shop    string `json:"shop"`
	OrderId int64  `json:"order_id"`
	AppId   string `json:"app_id"`
}

type OrderListReq struct {
	Uid      int    `json:"uid"`
	Page     int    `form:"page" binding:"required,gte=1"` // gte=1 表示必须 ≥1
	PageSize int    `form:"page_size" binding:"required,gte=1,lte=50"`
	Type     string `form:"type" binding:"required,oneof=All Paid Refund"`
	Query    string `form:"query" binding:"omitempty,max=20"`
}
