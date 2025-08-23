package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterPluginRouter(r *gin.RouterGroup, h *handler.SettingHandler) {
	// 对外
	publicGroup := r.Group("plugin")

	publicGroup.POST("/config", h.GetPublicCart)
}
