package jobs

import (
	"context"

	"backend/internal/domain/entity/jobs"
)

type OrderRepository interface {
	Create(ctx context.Context, jobOrder *jobs.JobOrder) (int, error)
	ExistsByOrderID(ctx context.Context, orderId string) int
	First(ctx context.Context, id int) (*jobs.JobOrder, error)
	UpdateJobTime(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status int) error
	Clear(ctx context.Context) error
}
