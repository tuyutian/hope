package jobs

import (
	"context"

	"backend/internal/domain/entity/jobs"
)

type OrderRepository interface {
	// First 查询订单任务
	First(ctx context.Context, jobId int64) (*jobs.JobOrder, error)
	// UpdateJobTime 更新任务时间
	UpdateJobTime(ctx context.Context, jobId int64) error
	// UpdateStatus 更新任务状态
	UpdateStatus(ctx context.Context, jobId int64, status int) error
}
