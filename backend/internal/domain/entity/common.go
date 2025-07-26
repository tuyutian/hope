package entity

type Pagination struct {
	Page int `json:"page" binding:"required,min=1"` // 页码
	Size int `json:"size" binding:"required,min=1"` // 每页数量
}
