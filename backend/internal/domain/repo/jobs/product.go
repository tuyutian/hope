package jobs

import (
	"context"

	"backend/internal/domain/entity/jobs"
)

type ProductRepository interface {
	First(ctx context.Context, id int64) (*jobs.JobProduct, error)
	Create(ctx context.Context, jobProduct *jobs.JobProduct) (int64, error)
	UpdateJobTime(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status int) error
	Clear(ctx context.Context) error
}
