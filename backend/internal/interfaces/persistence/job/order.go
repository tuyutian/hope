package job

import (
	"context"
	"time"

	"backend/internal/domain/entity/jobs"
	jobRepo "backend/internal/domain/repo/jobs"

	"xorm.io/xorm"
)

var _ jobRepo.OrderRepository = (*OrderRepoImpl)(nil)

type OrderRepoImpl struct {
	db *xorm.Engine
}

func NewOrderRepository(db *xorm.Engine) jobRepo.OrderRepository {
	return &OrderRepoImpl{db: db}
}

func (j *OrderRepoImpl) Create(ctx context.Context, jobOrder *jobs.JobOrder) (int64, error) {
	_, err := j.db.Context(ctx).Insert(jobOrder)
	if err != nil {
		return 0, err
	}
	return jobOrder.Id, nil
}

func (j *OrderRepoImpl) ExistsByOrderID(ctx context.Context, orderId int64) int64 {
	var jobOrder jobs.JobOrder
	has, err := j.db.Context(ctx).Cols("id").Where("order_id = ? and is_success = 0", orderId).Get(&jobOrder)

	if err != nil || !has {
		return 0
	}
	return 1
}

func (j *OrderRepoImpl) First(ctx context.Context, id int64) (*jobs.JobOrder, error) {
	var jobOrder jobs.JobOrder
	has, err := j.db.Context(ctx).Where("id = ?", id).Get(&jobOrder)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &jobOrder, nil
}

func (j *OrderRepoImpl) UpdateJobTime(ctx context.Context, id int64) error {
	_, err := j.db.Context(ctx).
		Table(new(jobs.JobOrder)).
		Where("id = ?", id).
		Update(map[string]interface{}{
			"job_time": time.Now().Unix(),
		})
	return err
}

func (j *OrderRepoImpl) UpdateStatus(ctx context.Context, id int64, status int) error {
	_, err := j.db.Context(ctx).
		Table(new(jobs.JobOrder)).
		Where("id = ?", id).
		Update(map[string]interface{}{
			"is_success": status,
		})
	return err
}

func (j *OrderRepoImpl) Clear(ctx context.Context) error {
	_, _ = j.db.Context(ctx).Where("is_success > ?", 0).Delete(&jobs.JobOrder{})
	return nil
}
