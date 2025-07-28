package orders

import "backend/internal/domain/entity"

type QueryOrderEntity struct {
	UserID int64 `json:"user_id,omitempty"`
	entity.Pagination
	Type  int    `json:"type" binding:"oneof=0 1 2"` // 订单类型
	Query string `json:"query"`
}

// OrderStatistics 用于存储查询结果
type OrderStatistics struct {
	TotalRefund     float64 `xorm:"'total_refund'" json:"total_refund"`
	TotalProtectify float64 `xorm:"'total_protectify'" json:"total_protectify"`
	TotalOrders     int     `xorm:"'total_orders'" json:"total_orders"`
}

type OrderWebHookReq struct {
	Shop    string `json:"shop"`
	OrderId int64  `json:"order_id"`
	AppId   string `json:"app_id"`
}

type OrderResponse struct {
	List  []*UserOrder `json:"list"`
	Total int64        `json:"total"`
}
