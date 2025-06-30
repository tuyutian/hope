package tasks

import (
	"github.com/hibiken/asynq"

	"backend/internal/interfaces/job/handler"
)

func InitTask(mux *asynq.ServeMux, handlers *handler.Handlers) {
	RegisterProductHandler(mux, handlers.ProductHandler)
	RegisterUserHandler(mux, handlers.UserHandler)
	RegisterOrderHandler(mux, handlers.OrderHandler)
}
