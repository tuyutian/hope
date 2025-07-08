package billings

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
)

type usageChargeGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

var _ shopifys.UsageChargeGraphqlRepository = (*usageChargeGraphqlRepoImpl)(nil)

func NewUsageChargeGraphqlRepository() shopifys.UsageChargeGraphqlRepository {
	return &usageChargeGraphqlRepoImpl{}
}

func (u *usageChargeGraphqlRepoImpl) CreateUsageCharge(ctx context.Context, lineItemId int64, amount decimal.Decimal, description string) (string, error) {

	// 构建 GraphQL 请求
	mutation := `
        mutation appUsageRecordCreate($description: String!, $price: MoneyInput!, $subscriptionLineItemId: ID!) {
            appUsageRecordCreate(description: $description, price: $price, subscriptionLineItemId: $subscriptionLineItemId) {
                appUsageRecord {
                    id
                    description
                    price {
                        amount
                        currencyCode
                    }
                    createdAt
                }
                userErrors {
                    field
                    message
                }
            }
        }
    `

	variables := map[string]interface{}{
		"description": description,
		"price": map[string]interface{}{
			"amount":       amount.String(),
			"currencyCode": "USD",
		},
		"subscriptionLineItemId": fmt.Sprintf("gid://shopify/SubscriptionLineItem/%d", lineItemId),
	}
	var response shopifyEntity.AppUsageRecordCreateResponse
	// 发送请求
	err := u.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return "", err
	}

	// 解析响应
	if len(response.UserErrors) > 0 {
		return "", errors.New(response.UserErrors[0].Message)
	}

	return response.AppUsageRecord.ID, nil
}
