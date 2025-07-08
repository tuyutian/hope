package billings

import (
	"context"
	"fmt"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
)

type subscriptionGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

var _ shopifys.SubscriptionGraphqlRepository = (*subscriptionGraphqlRepoImpl)(nil)

func NewSubscriptionGraphqlRepository() shopifys.SubscriptionGraphqlRepository {
	return &subscriptionGraphqlRepoImpl{}
}

// CreateSubscription 创建订阅
func (s *subscriptionGraphqlRepoImpl) CreateSubscription(ctx context.Context, input shopifyEntity.AppSubscriptionCreateInput) (*shopifyEntity.AppSubscription, string, error) {
	// GraphQL mutation
	mutation := `
        mutation appSubscriptionCreate($name: String!, $test: Boolean, $trialDays: Int, $lineItems: [AppSubscriptionLineItemInput!]!, $returnUrl: URL!) {
            appSubscriptionCreate(
                name: $name
                test: $test
                trialDays: $trialDays
                lineItems: $lineItems
                returnUrl: $returnUrl
            ) {
                appSubscription {
                    id
                    name
                    status
                    test
                    createdAt
                    currentPeriodEnd
                    trialDays
                    lineItems {
                        id
                        plan {
                            pricingDetails {
                                ... on AppUsagePricingDetails {
                                    cappedAmount {
                                        amount
                                        currencyCode
                                    }
                                    terms
                                    balanceUsed {
                                        amount
                                        currencyCode
                                    }
                                }
                                ... on AppRecurringPricingDetails {
                                    price {
                                        amount
                                        currencyCode
                                    }
                                    interval
                                }
                            }
                        }
                    }
                }
                confirmationUrl
                userErrors {
                    field
                    message
                }
            }
        }
    `

	// 构建变量
	variables := map[string]interface{}{
		"name":      input.Name,
		"test":      input.Test,
		"lineItems": input.LineItems,
		"returnUrl": input.ReturnURL,
	}

	if input.TrialDays > 0 {
		variables["trialDays"] = input.TrialDays
	}

	// 发送请求
	var response struct {
		AppSubscriptionCreate shopifyEntity.AppSubscriptionCreateResponse `json:"appSubscriptionCreate"`
	}

	err := s.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return nil, "", err
	}

	// 检查用户错误
	if len(response.AppSubscriptionCreate.UserErrors) > 0 {
		return nil, "", fmt.Errorf("shopify error: %s",
			response.AppSubscriptionCreate.UserErrors[0].Message)
	}

	return response.AppSubscriptionCreate.AppSubscription,
		response.AppSubscriptionCreate.ConfirmationURL, nil
}

// GetCurrentSubscription 获取当前订阅
func (s *subscriptionGraphqlRepoImpl) GetCurrentSubscription(ctx context.Context) (*shopifyEntity.AppSubscription, error) {
	query := `
        query getCurrentSubscription {
            currentAppInstallation {
                id
                activeSubscriptions {
                    id
                    name
                    status
                    test
                    createdAt
                    currentPeriodEnd
                    trialDays
                    lineItems {
                        id
                        plan {
                            pricingDetails {
                                ... on AppUsagePricingDetails {
                                    cappedAmount {
                                        amount
                                        currencyCode
                                    }
                                    terms
                                    balanceUsed {
                                        amount
                                        currencyCode
                                    }
                                }
                                ... on AppRecurringPricingDetails {
                                    price {
                                        amount
                                        currencyCode
                                    }
                                    interval
                                }
                            }
                        }
                    }
                }
            }
        }
    `

	var result struct {
		CurrentAppInstallation struct {
			ID                  string                          `json:"id"`
			ActiveSubscriptions []shopifyEntity.AppSubscription `json:"activeSubscriptions"`
		} `json:"currentAppInstallation"`
	}

	err := s.Client.Query(ctx, query, nil, &result)
	if err != nil {
		return nil, err
	}

	// 返回第一个活跃订阅
	if len(result.CurrentAppInstallation.ActiveSubscriptions) > 0 {
		return &result.CurrentAppInstallation.ActiveSubscriptions[0], nil
	}

	return nil, fmt.Errorf("no active subscription found")
}
