// Package shops
package shops

import (
	"context"
	"fmt"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
)

var _ shopifyRepo.ShopGraphqlRepository = (*shopGraphqlRepoImpl)(nil)

type shopGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

// NewShopGraphqlRepository 创建店铺GraphQL仓储实现
func NewShopGraphqlRepository() shopifyRepo.ShopGraphqlRepository {
	return &shopGraphqlRepoImpl{}
}

// GetShopInfo 获取店铺信息
func (c *shopGraphqlRepoImpl) GetShopInfo(ctx context.Context) (*shopifyEntity.Shop, error) {
	query := `
		query {
			shop {
				id
				name
				email
				myshopifyDomain
				primaryDomain {
					url
					host
				}
				plan {
					displayName
					partnerDevelopment
					shopifyPlus
				}
				url
				currencyCode
				currencyFormats {
					moneyFormat
					moneyInEmailsFormat
					moneyWithCurrencyFormat
					moneyWithCurrencyInEmailsFormat
				}
				enabledPresentmentCurrencies
				billingAddress {
					address1
					address2
					city
					zip
					provinceCode
					countryCodeV2
				}
				contactEmail
				customerAccounts
				ianaTimezone
				metafields(first: 10) {
					edges {
						node {
							id
							namespace
							key
							value
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{}

	var response shopifyEntity.ShopResponse
	err := c.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("查询店铺信息失败: %w", err)
	}

	return &response.Shop, nil
}

// UpdateShopBillingAddress 更新店铺账单地址
func (c *shopGraphqlRepoImpl) UpdateShopBillingAddress(ctx context.Context, input shopifyEntity.ShopBillingAddressInput) error {
	mutation := `
		mutation shopBillingAddressUpdate($billingAddress: MailingAddressInput!) {
			shopBillingAddressUpdate(billingAddress: $billingAddress) {
				shop {
					id
					billingAddress {
						address1
						address2
						city
						zip
						provinceCode
						countryCodeV2
					}
				}
				userErrors {
					field
					message
				}
			}
		}
	`

	variables := map[string]interface{}{
		"billingAddress": input,
	}

	var response struct {
		ShopBillingAddressUpdate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"shopBillingAddressUpdate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return fmt.Errorf("更新店铺账单地址失败: %w", err)
	}

	if len(response.ShopBillingAddressUpdate.UserErrors) > 0 {
		return fmt.Errorf("更新店铺账单地址错误: %s", response.ShopBillingAddressUpdate.UserErrors[0].Message)
	}

	return nil
}

// UpdateShopSettings 更新店铺设置
func (c *shopGraphqlRepoImpl) UpdateShopSettings(ctx context.Context, input shopifyEntity.ShopSettingsInput) error {
	mutation := `
		mutation shopSettingsUpdate($settings: ShopSettingsInput!) {
			shopSettingsUpdate(settings: $settings) {
				shop {
					id
					name
					contactEmail
				}
				userErrors {
					field
					message
				}
			}
		}
	`

	variables := map[string]interface{}{
		"settings": input,
	}

	var response struct {
		ShopSettingsUpdate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"shopSettingsUpdate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return fmt.Errorf("更新店铺设置失败: %w", err)
	}

	if len(response.ShopSettingsUpdate.UserErrors) > 0 {
		return fmt.Errorf("更新店铺设置错误: %s", response.ShopSettingsUpdate.UserErrors[0].Message)
	}

	return nil
}

// GetShopPolicies 获取店铺政策
func (c *shopGraphqlRepoImpl) GetShopPolicies(ctx context.Context) (*shopifyEntity.ShopPoliciesResponse, error) {
	query := `
		query {
			shop {
				id
				privacyPolicy {
					id
					title
					body
					url
				}
				refundPolicy {
					id
					title
					body
					url
				}
				termsOfService {
					id
					title
					body
					url
				}
				shippingPolicy {
					id
					title
					body
					url
				}
				subscriptionPolicy {
					id
					title
					body
					url
				}
			}
		}
	`

	variables := map[string]interface{}{}

	var response shopifyEntity.ShopPoliciesResponse
	err := c.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("查询店铺政策失败: %w", err)
	}

	return &response, nil
}

// GetShopLocales 获取店铺语言设置
func (c *shopGraphqlRepoImpl) GetShopLocales(ctx context.Context) (*shopifyEntity.ShopLocalesResponse, error) {
	query := `
		query {
			shop {
				id
				primaryDomain {
					url
				}
				primaryLocale
				url
				translatableContentV2(first: 10) {
					edges {
						node {
							key
							value
							digest
							translatableContent {
								key
								value
								digest
							}
							translations(first: 10) {
								locale
								key
								value
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{}

	var response shopifyEntity.ShopLocalesResponse
	err := c.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("查询店铺语言设置失败: %w", err)
	}

	return &response, nil
}

// CreateWebhookSubscription creates a new webhook subscription
func (c *shopGraphqlRepoImpl) CreateWebhookSubscription(ctx context.Context, topic string, callbackUrl string) error {
	mutation := `
		mutation webhookSubscriptionCreate($topic: WebhookSubscriptionTopic!, $webhookSubscription: WebhookSubscriptionInput!) {
			webhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription) {
				webhookSubscription {
					id
					topic
					filter
					format
					endpoint {
						__typename
						... on WebhookHttpEndpoint {
							callbackUrl
						}
					}
				}
				userErrors {
					field
					message
				}
			}
		}
	`

	variables := map[string]interface{}{
		"topic": topic,
		"webhookSubscription": map[string]interface{}{
			"callbackUrl": callbackUrl,
			"format":      "JSON",
		},
	}

	var response struct {
		WebhookSubscriptionCreate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"webhookSubscriptionCreate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return fmt.Errorf("创建webhook订阅失败: %w", err)
	}

	if len(response.WebhookSubscriptionCreate.UserErrors) > 0 {
		return fmt.Errorf("创建webhook订阅错误: %s", response.WebhookSubscriptionCreate.UserErrors[0].Message)
	}

	return nil
}

// UpdateWebhookSubscription updates an existing webhook subscription
func (c *shopGraphqlRepoImpl) UpdateWebhookSubscription(ctx context.Context, id string, callbackUrl string) error {
	mutation := `
		mutation webhookSubscriptionUpdate($id: ID!, $webhookSubscription: WebhookSubscriptionInput!) {
			webhookSubscriptionUpdate(id: $id, webhookSubscription: $webhookSubscription) {
				userErrors {
					field
					message
				}
				webhookSubscription {
					id
					topic
					endpoint {
						... on WebhookHttpEndpoint {
							callbackUrl
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
		"webhookSubscription": map[string]interface{}{
			"callbackUrl": callbackUrl,
			"format":      "JSON",
		},
	}

	var response struct {
		WebhookSubscriptionUpdate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"webhookSubscriptionUpdate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return fmt.Errorf("更新webhook订阅失败: %w", err)
	}

	if len(response.WebhookSubscriptionUpdate.UserErrors) > 0 {
		return fmt.Errorf("更新webhook订阅错误: %s", response.WebhookSubscriptionUpdate.UserErrors[0].Message)
	}

	return nil
}

// QueryWebhookSubscriptions gets a list of webhook subscriptions
func (c *shopGraphqlRepoImpl) QueryWebhookSubscriptions(ctx context.Context, queryParams string) ([]shopifyEntity.WebhookSubscription, error) {
	query := `
		query ($query: String) {
			webhookSubscriptions(first: 20, query: $query) {
				nodes {
					apiVersion {
						displayName
						handle
						supported
					}
					createdAt
					endpoint {
						... on WebhookHttpEndpoint {
							__typename
							callbackUrl
						}
					}
					format
					filter
					id
					includeFields
					metafieldNamespaces
					topic
					updatedAt
				}
			}
		}
	`

	variables := map[string]interface{}{
		"query": queryParams,
	}

	var response struct {
		WebhookSubscriptions struct {
			Nodes []shopifyEntity.WebhookSubscription `json:"nodes"`
		} `json:"webhookSubscriptions"`
	}

	err := c.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("查询webhook订阅列表失败: %w", err)
	}

	return response.WebhookSubscriptions.Nodes, nil
}

func (c *shopGraphqlRepoImpl) GetPublicationList(ctx context.Context) (string, error) {
	query := `
		query {
			publications(first: 20) {
				edges {
					node {
						id
						name
					}
				}
			}
		}
		`
	// 初始化返回的数据结构
	var response struct {
		Publications struct {
			Edges []struct {
				Node struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"publications"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err := c.Client.Query(ctx, query, nil, &response)
	if err != nil {
		return "", err
	}

	if len(response.Errors) > 0 {
		return "", fmt.Errorf("shopify 返回错误: %+v", response.Errors)
	}

	// 查找 "Online Store"
	for _, edge := range response.Publications.Edges {
		if edge.Node.Name == "Online Store" {
			return edge.Node.ID, nil
		}
	}
	// 返回查询结果
	return "", nil
}
