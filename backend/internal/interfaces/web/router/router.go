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

// InitRouters 初始化router规则
func InitRouters(router *gin.Engine, handlers *handler.Handlers) {
	// 访问日志中间件处理
	logWare := &middleware.LogWare{}

	// 对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(
		logWare.Access(), middleware.Cors(), func() gin.HandlerFunc {
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
	router.NoRoute(middleware.NotFoundHandler())

}
