package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorsWare cors middleware
type CorsWare struct {
	AllowedOrigins []string
}

func (ware *CorsWare) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		ware.setCorsHeaders(c)

		// 处理OPTIONS预检请求Add commentMore actions
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

var (
	allowMethod  = "POST,GET,OPTIONS,PUT,DELETE,UPDATE"
	allowHeaders = "Content-Type,AccessToken,Authorization,Cookie,cookie,Content-Length,X-CSRF-Token," +
		"Token,session,ignorecanceltoken,X-Requested-With"
	exposeHeaders = "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type"
)

// setCorsHeaders 设置CORS相关的响应头
func (ware *CorsWare) setCorsHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", allowMethod)
	c.Header("Access-Control-Allow-Headers", allowHeaders)
	c.Header("Access-Control-Expose-Headers", exposeHeaders)
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "true")
}
