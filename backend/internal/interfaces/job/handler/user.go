package handler

import (
	"context"

	"github.com/hibiken/asynq"

	"backend/internal/application/jobs"
)

type UserHandler struct {
	userService *jobs.UserService
}

func (h *UserHandler) HandleInitUser(ctx context.Context, task *asynq.Task) error {
	return h.userService.HandleInitUser(ctx, task)
}
