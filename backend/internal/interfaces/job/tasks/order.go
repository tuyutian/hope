package tasks

import (
	"github.com/hibiken/asynq"

	"backend/internal/infras/config"
	"backend/internal/interfaces/job/handler"
)

func RegisterOrderHandler(mux *asynq.ServeMux, handler *handler.OrderHandler) {

	mux.HandleFunc(config.SendOrder, handler.HandleOrder)
	mux.HandleFunc(config.SendOrderStatistics, handler.HandleOrderStatistics)

}
