package shopifys

import (
	"context"

	shopifyEntity "backend/internal/domain/entity/shopifys"
)

type OrderGraphqlRepository interface {
	BaseGraphqlRepository

	GetOrderInfo(ctx context.Context, orderId int64) (*shopifyEntity.OrderResponse, error)
}
