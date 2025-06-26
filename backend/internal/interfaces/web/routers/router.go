package routers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
	"backend/pkg/utils"
)

// Middleware 路由中间件
type Middleware struct {
	AuthWare    *middleware.AuthWare
	RequestWare *middleware.RequestWare
	CorsWare    *middleware.CorsWare
}

// InitRouters 初始化router规则
func InitRouters(router *gin.Engine, handlers *handler.Handlers, middlewares *Middleware) {
	requestWare := middlewares.RequestWare

	// 对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(
		requestWare.Access(), middlewares.CorsWare.Cors(), func() gin.HandlerFunc {
			return func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("Panic info is: %v", err)
						go utils.CallWilding(fmt.Sprintf("Painc App:%s\nInfo is: %v\nRequest is: %s\nDomain: %s", "tms-api", err, c.Request.URL, c.GetHeader("X-Shopify-Shop-Domain")))
					}
				}()
			}
		}(), middleware.TimeoutHandler(40*time.Second),
	)

	// gin 框架prometheus接入
	//router.Use(middleware.WrapMonitor())

	// 路由找不到的情况
	router.NoRoute(requestWare.NotFoundHandler())
	api := router.Group("api/v1") // 定义路由组
	RegisterOrderRouter(api, handlers.OrderHandler, middlewares.AuthWare)
}
