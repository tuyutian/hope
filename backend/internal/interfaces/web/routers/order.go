package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterOrderRouter(r *gin.RouterGroup, h *handler.OrderHandler, m *Middleware) {
	orderGroup := r.Group("/order", m.AuthWare.CheckLogin(), m.ShopifyGraphqlWare.ShopifyGraphqlClient())

	orderGroup.GET("/orders", h.OrderList)
}
