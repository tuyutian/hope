package shopifys

import (
	"context"

	"backend/internal/infras/shopify_graphql"
)

type ShopifyRepository interface {
	ExtractCurrencySymbol(moneyFormat string) string
	GetWebhookUrl(appID string) string
	GetReturnUrl(appID string, userID int64) string
	VerifyWebhook(ctx context.Context, appSecret string, signature string, body []byte) bool
}

type BaseGraphqlRepository interface {
	WithClient(client *shopify_graphql.GraphqlClient)
}
