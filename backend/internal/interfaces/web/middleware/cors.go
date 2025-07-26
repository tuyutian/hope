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
		origin := c.Request.Header.Get("Origin")
		ware.setCorsHeaders(c, origin)

		if origin != "" {
			// 接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//
			// c.Header("Content-Type", "application/json")
		}

		// 允许类型校验
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
		"Token,session,ignorecanceltoken"
	exposeHeaders = "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type"
)

// setCorsHeaders 设置CORS相关的响应头
func (ware *CorsWare) setCorsHeaders(c *gin.Context, origin string) {
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", allowMethod)
	c.Header("Access-Control-Allow-Headers", allowHeaders)
	c.Header("Access-Control-Expose-Headers", exposeHeaders)
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "true")
}
