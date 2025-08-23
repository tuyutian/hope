package middleware

import (
	"net/http"
	"strings"

	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CorsWare cors middleware
type CorsWare struct {
	AllowedOrigins []string
}

func (ware *CorsWare) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 如果是直接访问（没有 Origin 头），允许通过
		if origin == "" {
			c.Next()
			return
		}

		// 检查是否允许该域名
		if ware.isAllowedOrigin(origin) {
			ware.setCorsHeaders(c, origin)
			logger.Info(c.Request.Context(), "CORS 允许", "origin", origin)
		} else {
			logger.Warn(c.Request.Context(), "CORS 拒绝", "origin", origin)
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查是否允许该域名
func (ware *CorsWare) isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	if strings.Contains(origin, "shopify.com") {
		return true
	}
	// 如果配置了允许的域名列表，使用配置的列表
	if len(ware.AllowedOrigins) > 0 {
		for _, allowed := range ware.AllowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	}

	// 默认允许的域名列表
	allowedOrigins := []string{
		"https://s.protectifyapp.com",
		"https://protectifyapp.com",
		"http://localhost:9527",
		"http://127.0.0.1:9527",
	}

	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

var (
	allowMethod  = "POST,GET,OPTIONS,PUT,DELETE,UPDATE"
	allowHeaders = "Content-Type,AccessToken,Authorization,Cookie,cookie,Content-Length,X-CSRF-Token," +
		"Token,session,ignorecanceltoken,X-Requested-With"
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
