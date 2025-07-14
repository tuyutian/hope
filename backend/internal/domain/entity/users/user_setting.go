package users

type UpdateStep struct {
	Name string `json:"name" binding:"required"`
	Open bool   `json:"open" binding:"required"`
}

type UpdateSetting struct {
	Name  string `json:"name" binding:"required"`                // 必填字段
	Value string `json:"value" binding:"required,min=1,max=100"` // 必填，且长度在1-100之间
}
