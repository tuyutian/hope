package billings

import (
	"context"

	"backend/internal/domain/entity"
	"backend/internal/domain/entity/billings"
)

type BillingPeriodSummaryRepository interface {
	BillingPeriodSummary(ctx context.Context, userID int64, pagination entity.Pagination) ([]*billings.BillingPeriodSummary, error)
	BillingPeriodCount(ctx context.Context, userID int64) (int64, error)
	CreateBillingPeriodSummary(ctx context.Context, period *billings.BillingPeriodSummary) (int64, error)
	UpdateBillingPeriodSummary(ctx context.Context, period *billings.BillingPeriodSummary) error
	GetByCurrentPeriod(ctx context.Context, userID int64, periodEnd int64) (*billings.BillingPeriodSummary, error)
}
