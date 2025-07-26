package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterSettingRouter(r *gin.RouterGroup, h *handler.SettingHandler, m *Middleware) {

	settingGroup := r.Group("setting")
	// 私有路由使用jwt验证
	settingGroup.Use(m.AuthWare.CheckLogin(), m.ShopifyGraphqlWare.ShopifyGraphqlClient())

	settingGroup.GET("/cart", h.GetCart)
	settingGroup.POST("/cart", h.UpdateCart)
}
