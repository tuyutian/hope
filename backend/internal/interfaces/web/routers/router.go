package routers

import (
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
)

// Middleware 路由中间件
type Middleware struct {
	AuthWare           *middleware.AuthWare
	RequestWare        *middleware.RequestWare
	CorsWare           *middleware.CorsWare
	CspWare            *middleware.CspMiddleware
	ShopifyGraphqlWare *middleware.ShopifyGraphqlWare
	AppMiddleware      *middleware.AppMiddleware
}

// InitRouters 初始化router规则
func InitRouters(router *gin.Engine, handlers *handler.Handlers, middlewares *Middleware) {
	requestWare := middlewares.RequestWare

	// 对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(
		requestWare.Access(),
		middlewares.CorsWare.Cors(),
		requestWare.Recover(),
		middleware.TimeoutHandler(40*time.Second),
	)

	// gin 框架prometheus接入
	//router.Use(middleware.WrapMonitor())

	// 路由找不到的情况
	router.NoRoute(requestWare.NotFoundHandler())
	api := router.Group("/:appId/api/v1") // 定义路由组
	api.Use(middlewares.AppMiddleware.AppMust(), middlewares.CspWare.Csp())
	RegisterPluginRouter(api, handlers.SettingHandler)
	RegisterWebhookRouter(api, handlers.WebhookHandler)
	RegisterCommonRouter(api, handlers.CommonHandler, middlewares.AuthWare)
	RegisterBillingRouter(api, handlers.BillingHandler, middlewares.AuthWare)
	RegisterSettingRouter(api, handlers.SettingHandler, middlewares)
	RegisterOrderRouter(api, handlers.OrderHandler, middlewares)
	RegisterUserRouter(api, handlers.UserHandler, middlewares)
}
