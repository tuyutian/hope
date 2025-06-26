package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
)

func RegisterOrderRouter(r *gin.RouterGroup, h *handler.OrderHandler, m *middleware.AuthWare) {
	orderGroup := r.Group("/order", m.CheckLogin())

	orderGroup.GET("/orders", h.OrderList)
}
