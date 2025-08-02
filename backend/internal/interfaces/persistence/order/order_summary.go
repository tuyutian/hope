package order

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity/orders"
	orderRepo "backend/internal/domain/repo/orders"
)

var _ orderRepo.OrderSummaryRepository = (*summaryRepoImpl)(nil)

type summaryRepoImpl struct {
	db *xorm.Engine
}

// NewOrderSummaryRepository NewSummaryRepository 从数据库获取订单统计资源
func NewOrderSummaryRepository(engine *xorm.Engine) orderRepo.OrderSummaryRepository {
	return &summaryRepoImpl{db: engine}
}

func (s *summaryRepoImpl) GetByDays(ctx context.Context, userId int64, days int) ([]orders.OrderSummary, error) {
	var summary []orders.OrderSummary
	err := s.db.
		Where("user_id = ? ", userId).
		Desc("id").
		Limit(days).
		Find(&summary)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *summaryRepoImpl) ExistOrder(ctx context.Context, userID int64, today int64) (int64, error) {
	var orderSummary orders.OrderSummary
	has, err := s.db.Context(ctx).
		Where("user_id = ? AND today = ?", userID, today).
		Cols("id").
		Get(&orderSummary)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, nil
	}

	return orderSummary.Id, nil
}

func (s *summaryRepoImpl) UpsertOrderStatistics(ctx context.Context, orderSummary orders.OrderSummary) error {
	_, err := s.db.Context(ctx).ID(orderSummary.Id).Update(&orderSummary)
	return err
}

func (s *summaryRepoImpl) CrateOrderStatistics(ctx context.Context, orderSummary orders.OrderSummary) error {
	_, err := s.db.Context(ctx).Insert(&orderSummary)
	return err
}
