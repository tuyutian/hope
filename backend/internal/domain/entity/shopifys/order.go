package shopifys

type Order struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	CreatedAt              string `json:"createdAt"`
	ProcessedAt            string `json:"processedAt"`
	DisplayFinancialStatus string `json:"displayFinancialStatus"`
	TotalPriceSet          struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"shopMoney"`
	} `json:"totalPriceSet"`
	LineItems struct {
		Edges []struct {
			Node struct {
				VariantTitle string `json:"variantTitle"`
				Sku          string `json:"sku"`
				Quantity     int    `json:"quantity"`
				Variant      struct {
					ID string `json:"id"`
				} `json:"variant"`
				OriginalUnitPriceSet struct {
					ShopMoney struct {
						Amount       string `json:"amount"`
						CurrencyCode string `json:"currencyCode"`
					} `json:"shopMoney"`
				} `json:"originalUnitPriceSet"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"lineItems"`
	Refunds []struct {
		ID               string `json:"id"`
		CreatedAt        string `json:"createdAt"`
		TotalRefundedSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currencyCode"`
			} `json:"shopMoney"`
		} `json:"totalRefundedSet"`
		RefundLineItems struct {
			Edges []struct {
				Node struct {
					LineItem struct {
						ID      string `json:"id"`
						Title   string `json:"title"`
						Variant struct {
							ID string `json:"id"`
						} `json:"variant"`
					} `json:"lineItem"`
					Quantity    int `json:"quantity"`
					SubtotalSet struct {
						ShopMoney struct {
							Amount       string `json:"amount"`
							CurrencyCode string `json:"currencyCode"`
						} `json:"shopMoney"`
					} `json:"subtotalSet"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"refundLineItems"`
	} `json:"refunds"`
}

type OrderResponse struct {
	Order Order `json:"order"`
}
