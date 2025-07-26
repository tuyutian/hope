package billing

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity"
	billingEntity "backend/internal/domain/entity/billings"
	"backend/internal/domain/repo/billings"
)

var _ billings.BillingPeriodSummaryRepository = (*billingPeriodSummaryRepoImpl)(nil)

type billingPeriodSummaryRepoImpl struct {
	db *xorm.Engine
}

func NewBillingPeriodSummaryRepo(db *xorm.Engine) billings.BillingPeriodSummaryRepository {
	return &billingPeriodSummaryRepoImpl{db: db}
}

func (b *billingPeriodSummaryRepoImpl) BillingPeriodSummary(ctx context.Context, userID int64, pagination entity.Pagination) ([]*billingEntity.BillingPeriodSummary, error) {
	var periods []*billingEntity.BillingPeriodSummary
	err := b.db.Table(new(billingEntity.BillingPeriodSummary)).Where("user_id = ?", userID).Desc("create_time").Limit(pagination.Size, (pagination.Page-1)*pagination.Size).Find(&periods)

	return periods, err
}

func (b *billingPeriodSummaryRepoImpl) BillingPeriodCount(ctx context.Context, userID int64) (int64, error) {
	var period billingEntity.BillingPeriodSummary
	count, err := b.db.Table(&period).Where("user_id = ?", userID).Count(period)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (b *billingPeriodSummaryRepoImpl) CreateBillingPeriodSummary(ctx context.Context, period *billingEntity.BillingPeriodSummary) (int64, error) {
	_, err := b.db.Table(new(billingEntity.BillingPeriodSummary)).Insert(period)
	if err != nil {
		return 0, err
	}
	return period.Id, nil
}

func (b *billingPeriodSummaryRepoImpl) UpdateBillingPeriodSummary(ctx context.Context, period *billingEntity.BillingPeriodSummary) error {
	_, err := b.db.Table(new(billingEntity.BillingPeriodSummary)).ID(period.Id).Update(period)
	if err != nil {
		return err
	}
	return nil
}

func (b *billingPeriodSummaryRepoImpl) GetByCurrentPeriod(ctx context.Context, userID int64, periodEnd int64) (*billingEntity.BillingPeriodSummary, error) {
	var period billingEntity.BillingPeriodSummary
	has, err := b.db.Table(&period).Where("user_id = ? and period_end = ?", userID, periodEnd).Get(&period)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &period, nil
}
