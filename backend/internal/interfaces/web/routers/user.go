package routers

import (
	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
)

func RegisterUserRouter(r *gin.RouterGroup, handler *handler.UserHandler, m *Middleware) {
	userGroup := r.Group("user")
	// 私有路由使用jwt验证
	userGroup.Use(m.AuthWare.CheckLogin(), m.ShopifyGraphqlWare.ShopifyGraphqlClient())
	userGroup.POST("step", handler.SetUserStep)
	userGroup.GET("conf", handler.GetUserConf)
	userGroup.GET("session", handler.GetSessionData)
	userGroup.POST("setting", handler.UpdateUserSetting)

}
