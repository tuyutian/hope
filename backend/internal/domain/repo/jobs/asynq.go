package jobs

import (
	"context"

	"github.com/hibiken/asynq"
)

type AsynqRepository interface {
	NewProductTask(ctx context.Context, jobId int64, userProductId int64, shopifyProductId int64) (*asynq.TaskInfo, error)
	InitWebhookUserTask(ctx context.Context, userID int64) (*asynq.TaskInfo, error)
	OrderWebhookTask(ctx context.Context, jobId int64) (*asynq.TaskInfo, error)
	ProductWebhookUpdateTask(ctx context.Context, userID int64, userProductId int64) (*asynq.TaskInfo, error)
	OrderStatisticsTask(ctx context.Context, userID int64, start int64, end int64) (*asynq.TaskInfo, error)
	DelProductTask(ctx context.Context, userID int64, productId int64, delType int) (*asynq.TaskInfo, error)
}
