package handler

import (
	"context"

	"github.com/hibiken/asynq"

	"backend/internal/application/jobs"
)

type OrderHandler struct {
	orderService *jobs.OrderService
}

func (h OrderHandler) HandleOrder(ctx context.Context, task *asynq.Task) error {
	return h.orderService.HandleOrder(ctx, task)
}

func (h OrderHandler) HandleOrderStatistics(ctx context.Context, task *asynq.Task) error {
	return h.orderService.HandleOrderStatistics(ctx, task)
}
