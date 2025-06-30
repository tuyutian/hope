package orders

import (
	orderEntity "backend/internal/domain/entity/orders"

	"context"
)

type OrderRepository interface {
	// DelOrder 软删除订单
	DelOrder(ctx context.Context, userID int64, orderId string) error
	// List 分页查询订单列表
	List(ctx context.Context, req orderEntity.QueryOrderEntity) ([]*orderEntity.UserOrder, int64, error)
	// Create 创建订单
	Create(ctx context.Context, order *orderEntity.UserOrder) (int64, error)
	// ExistsByOrderID 检查订单是否存在
	ExistsByOrderID(ctx context.Context, orderId int64, userID int64) int64
	// UpdateShopifyOrderId 更新订单信息
	UpdateShopifyOrderId(ctx context.Context, order *orderEntity.UserOrder) error
	// GetOrderStatistics 获取订单统计信息
	GetOrderStatistics(ctx context.Context, start, end int64, userID int64) (*orderEntity.OrderStatistics, error)
}
