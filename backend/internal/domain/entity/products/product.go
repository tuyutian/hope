package products

type ProductReq struct {
	UserID      int64  `json:"user_id"`     // 用户ID
	ProductType string `json:"ProductType"` // 产品type
	Collection  string `json:"Collection"`  // 产品集合
}

type ProductWebHookReq struct {
	Shop      string `json:"shop"`
	AppId     string `json:"app_id"`
	ProductId int64  `json:"product_id"`
}
