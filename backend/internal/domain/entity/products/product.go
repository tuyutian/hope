package products

type ProductWebHookReq struct {
	Shop      string `json:"shop"`
	AppId     string `json:"app_id"`
	ProductId int64  `json:"product_id"`
}
