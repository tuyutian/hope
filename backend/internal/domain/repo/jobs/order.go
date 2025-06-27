package jobs

import (
	"context"

	"backend/internal/domain/entity/jobs"
)

type OrderRepository interface {
	Create(ctx context.Context, jobOrder *jobs.JobOrder) (int64, error)
	ExistsByOrderID(ctx context.Context, orderId string) int64
	First(ctx context.Context, id int64) (*jobs.JobOrder, error)
	UpdateJobTime(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status int) error
	Clear(ctx context.Context) error
}
