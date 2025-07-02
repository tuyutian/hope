package shopifys

import "backend/internal/infras/shopify_graphql"

type ShopifyRepository interface {
	ExtractCurrencySymbol(moneyFormat string) string
	GetWebhookUrl() string
}

type BaseGraphqlRepository interface {
	WithClient(client *shopify_graphql.GraphqlClient)
}
