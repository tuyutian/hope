package orders

type QueryOrderEntity struct {
	UserID   int    `json:"user_id"`
	Page     int    `form:"page" binding:"required,gte=1"` // gte=1 表示必须 ≥1
	PageSize int    `form:"page_size" binding:"required,gte=1,lte=50"`
	Type     string `form:"type" binding:"required,oneof=All Paid Refund"`
	Query    string `form:"query" binding:"omitempty,max=20"`
}
