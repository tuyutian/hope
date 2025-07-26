package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterWebhookRouter(r *gin.RouterGroup, handler *handler.WebHookHandler) {
	//publicRouter := parent.Group("category")
	webhookGroup := r.Group("webhook")
	/*webhookGroup.POST("/uninstall", handler.Uninstall)
	webhookGroup.POST("/shop/update", handler.ShopUpdate)

	webhookGroup.POST("/order-fulfilled", handler.OrderUpdate)
	webhookGroup.POST("/order-paid", handler.OrderUpdate)
	webhookGroup.POST("/order-partially_fulfilled", handler.OrderUpdate)
	webhookGroup.POST("/order-updated", handler.OrderUpdate)
	webhookGroup.POST("/order-delete", handler.OrderDel)

	webhookGroup.POST("/product-update", handler.ProductUpdate)
	webhookGroup.POST("/product-delete", handler.ProductDel)

	webhookGroup.POST("/customer-data_request", handler.Customers)
	webhookGroup.POST("/customer-redact", handler.Customers)*/
	webhookGroup.POST("/shopify", handler.Shopify)
	webhookGroup.GET("/charge_callback/:userID", handler.ChargeCallback)
}
