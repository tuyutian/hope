package tasks

import (
	"github.com/hibiken/asynq"

	"backend/internal/infras/config"
	"backend/internal/interfaces/job/handler"
)

func RegisterProductHandler(mux *asynq.ServeMux, handler *handler.ProductHandler) {
	mux.HandleFunc(config.SendProduct, handler.HandleProduct)
	mux.HandleFunc(config.SendDelProduct, handler.HandleDelProduct)
	mux.HandleFunc(config.SendUpdateProduct, handler.HandleShopifyProduct)

}
