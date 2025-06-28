package apps

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"

	appEntity "backend/internal/domain/entity/apps"
)

type AppRepository interface {
	GetAppDefinition(ctx context.Context, appId string) (*appEntity.AppDefinition, error)
	GetGoShopifyByAppID(ctx context.Context, appId string) (goshopify.App, error)
}
