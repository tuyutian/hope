package shopifys

import (
	"context"

	"backend/internal/domain/entity/shopifys"
)

var (
	ShopifyWebhookTopics = []string{
		"orders/updated",
		"orders/delete",
		"app_subscriptions/update",
		"app_subscriptions/cancel",
		"app_subscriptions/approaching_capped_amount",
		"app/uninstalled",
		"products/update",
		"products/delete",
		"shop/update",
	}
	ShopifyComplianceTopics = []string{
		"customers/data_request", "customers/redact", "shop/redact",
	}
)

// ShopGraphqlRepository 店铺GraphQL仓储接口
type ShopGraphqlRepository interface {
	BaseGraphqlRepository
	GetShopInfo(ctx context.Context) (*shopifys.Shop, *shopifys.CurrentAppInstallation, error)
	UpdateShopBillingAddress(ctx context.Context, input shopifys.ShopBillingAddressInput) error
	UpdateShopSettings(ctx context.Context, input shopifys.ShopSettingsInput) error
	GetShopPolicies(ctx context.Context) (*shopifys.ShopPoliciesResponse, error)
	GetShopLocales(ctx context.Context) (*shopifys.ShopLocalesResponse, error)
	GetPublicationID(ctx context.Context) (string, error)
	QueryWebhookSubscriptions(ctx context.Context, queryParams string) ([]shopifys.WebhookSubscription, error)
	CreateWebhookSubscription(ctx context.Context, topic string, callbackUrl string) error
	UpdateWebhookSubscription(ctx context.Context, id string, callbackUrl string) error
	MetafieldSet(ctx context.Context, ownerId string, namespace string, fieldType string, key string, value string) (*[]shopifys.Metafield, error)
}

type ThemeGraphqlRepository interface {
	GetMainThemeSettingJson(ctx context.Context) (string, error)
}
