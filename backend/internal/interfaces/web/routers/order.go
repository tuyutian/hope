package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterOrderRouter(r *gin.RouterGroup, h *handler.OrderHandler, m *Middleware) {
	// 公开的测试接口，不需要认证
	r.GET("/order/test-dashboard", h.TestDashboard)

	orderGroup := r.Group("/order", m.AuthWare.CheckLogin(), m.ShopifyGraphqlWare.ShopifyGraphqlClient())

	orderGroup.POST("/list", h.OrderList)
	orderGroup.GET("/dashboard", h.Dashboard)
}
