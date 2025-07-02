package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterWebhookRouter(r *gin.RouterGroup, handler *handler.WebHookHandler) {
	//publicRouter := parent.Group("category")
	privateRouter := r.Group("webhook/:appId")
	privateRouter.POST("/uninstall", handler.Uninstall)
	privateRouter.POST("/shop/update", handler.ShopUpdate)

	privateRouter.POST("/order-fulfilled", handler.OrderUpdate)
	privateRouter.POST("/order-paid", handler.OrderUpdate)
	privateRouter.POST("/order-partially_fulfilled", handler.OrderUpdate)
	privateRouter.POST("/order-updated", handler.OrderUpdate)
	privateRouter.POST("/order-delete", handler.OrderDel)

	privateRouter.POST("/product-update", handler.ProductUpdate)
	privateRouter.POST("/product-delete", handler.ProductDel)

	privateRouter.POST("/customer-data_request", handler.Customers)
	privateRouter.POST("/customer-redact", handler.Customers)
	privateRouter.POST("/shop-redact", handler.Customers)
}
