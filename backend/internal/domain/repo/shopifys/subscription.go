package shopifys

import (
	"context"

	"github.com/shopspring/decimal"

	shopifyEntity "backend/internal/domain/entity/shopifys"
)

type SubscriptionGraphqlRepository interface {
	BaseGraphqlRepository
	CreateSubscription(ctx context.Context, input shopifyEntity.AppSubscriptionCreateInput) (*shopifyEntity.AppSubscription, string, error)
	GetCurrentSubscription(ctx context.Context) (*shopifyEntity.AppSubscription, error)
	GetRecurrentChargeByID(ctx context.Context, id int64) (*shopifyEntity.AppSubscription, error)
}

type UsageChargeGraphqlRepository interface {
	BaseGraphqlRepository
	CreateUsageCharge(ctx context.Context, lineItemId string, amount decimal.Decimal, description string) (string, error)
}
