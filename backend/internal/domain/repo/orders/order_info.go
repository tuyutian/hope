package orders

import (
	"backend/internal/domain/entity/orders"

	"context"
)

type OrderInfoRepository interface {
	// Create 创建订单详情
	Create(ctx context.Context, orderInfo []*orders.UserOrderInfo) error
	// UpdateShopifyVariants 更新Shopify变体
	UpdateShopifyVariants(ctx context.Context, userOrderId int64, variantId string, orderInfo *orders.UserOrderInfo) error
	// GetOrderDetailVariantIDs 获取订单详情变体ID列表
	GetOrderDetailVariantIDs(ctx context.Context, userOrderId int64, userID int64) ([]string, error)
}
