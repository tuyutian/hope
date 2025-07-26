package jobs

type ProductPayload struct {
	JobId            int64 `json:"job_id"`
	UserProductId    int64 `json:"user_product_id"`
	ShopifyProductId int64 `json:"shopify_product_id"`
}

type InitUserPayload struct {
	UserID int64 `json:"user_id"`
}

type OrderPayload struct {
	JobId int64 `json:"job_id"`
}

type ShopifyProductPayload struct {
	UserID        int64 `json:"user_id"`
	UserProductId int64 `json:"user_product_id"`
}

type OrderStatisticPayload struct {
	UserID int64 `json:"user_id"`
	Start  int64 `json:"start"`
	End    int64 `json:"end"`
}

type DelProductPayload struct {
	UserID    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
	DelType   int   `json:"del_type"`
}
