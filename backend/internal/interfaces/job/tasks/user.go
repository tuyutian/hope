package tasks

import (
	"github.com/hibiken/asynq"

	"backend/internal/infras/config"
	"backend/internal/interfaces/job/handler"
)

func RegisterUserHandler(mux *asynq.ServeMux, handler *handler.UserHandler) {
	mux.HandleFunc(config.SendInitUser, handler.HandleInitUser)

}
