package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
)

type CommonRouter struct{}

func RegisterCommonRouter(r *gin.RouterGroup, h *handler.CommonHandler, authWare *middleware.AuthWare) {
	//publicRouter := parent.Group("category")
	commonGroup := r.Group("common")
	// 私有路由使用jwt验证
	commonGroup.Use(authWare.CheckLogin(), authWare.CheckAdmin())
	// 依赖注入
	commonGroup.POST("upload", h.Upload)

}
