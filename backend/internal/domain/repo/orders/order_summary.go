package orders

import (
	"context"

	"backend/internal/domain/entity/orders"
)

type OrderSummaryRepository interface {
	// GetByDays FindByToday 查询今日订单统计
	GetByDays(ctx context.Context, userId int64, days int) ([]orders.OrderSummary, error)
	// ExistOrder 检查是否存在指定日期的订单统计
	ExistOrder(ctx context.Context, userID int64, today int64) (int64, error)
	// UpsertOrderStatistics 更新订单统计
	UpsertOrderStatistics(ctx context.Context, orderSummary orders.OrderSummary) error
	// CrateOrderStatistics 创建订单统计
	CrateOrderStatistics(ctx context.Context, orderSummary orders.OrderSummary) error
}
