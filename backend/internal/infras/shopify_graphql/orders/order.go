package orders

import (
	"context"
	"fmt"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
)

var _ shopifys.OrderGraphqlRepository = (*orderGraphqlRepoImpl)(nil)

type orderGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

func NewOrderGraphqlRepository() shopifys.OrderGraphqlRepository {
	return &orderGraphqlRepoImpl{}
}

func (o *orderGraphqlRepoImpl) GetOrderInfo(ctx context.Context, orderId int64) (*shopifyEntity.OrderResponse, error) {
	orderGId := fmt.Sprintf("gid://shopify/Order/%d", orderId)
	query := `
		query($id: ID!) {
		  order(id: $id) {
			id
			name
			email
			displayFinancialStatus
			processedAt
			createdAt
			totalPriceSet {
			  shopMoney {
				amount
				currencyCode
			  }
			}
			lineItems(first: 100) {
			  edges {
				node {
			      variantTitle
				  sku
				  quantity
				  variant {
					id
				  }
				  originalUnitPriceSet {
					shopMoney {
					  amount
					  currencyCode
					}
				  }
				}
			  }
			}
			refunds {
			  id
			  createdAt
			  totalRefundedSet {
				shopMoney {
				  amount
				  currencyCode
				}
			  }
			  refundLineItems(first: 100) {
				edges {
				  node {
					lineItem {
					  id
					  variant {
						id
					  }
					}
					quantity
					subtotalSet {
					  shopMoney {
						amount
						currencyCode
					  }
					}
				  }
				}
			  }
			}
		  }
		}
	`

	var response shopifyEntity.OrderResponse
	vars := map[string]interface{}{
		"id": orderGId,
	}
	err := o.Client.Query(ctx, query, vars, &response)
	if err != nil {
		return nil, err
	}

	if response.Order.ID == "" {
		return nil, fmt.Errorf("order not found")
	}

	return &response, nil
}
