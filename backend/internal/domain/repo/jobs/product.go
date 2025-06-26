package jobs

import (
	"context"

	"backend/internal/domain/entity/jobs"
)

type ProductRepository interface {
	First(ctx context.Context, id int) (*jobs.JobProduct, error)
	Create(ctx context.Context, jobProduct *jobs.JobProduct) (int, error)
	UpdateJobTime(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status int) error
	Clear(ctx context.Context) error
}
