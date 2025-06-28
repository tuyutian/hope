package shopifys

import (
	"context"

	shopifyEntity "backend/internal/domain/entity/shopifys"
)

type ShopifyRepository interface {
	RequestOfflineSessionToken(ctx context.Context, token string) (*shopifyEntity.Token, error)
	GetShopName(ctx context.Context, shopUrl string) (string, error)
	ExtractCurrencySymbol(moneyFormat string) string
}
