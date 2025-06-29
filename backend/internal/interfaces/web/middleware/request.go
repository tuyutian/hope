package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/pkg/logger"
)

// NotFoundHandler router not found.
func NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"message": "this page not found!",
		})
	}
}

// TimeoutHandler server timeout middleware wraps the request context with a timeout
// 中间件参考go-chi/chi库 https://github.com/go-chi/chi/blob/master/middleware/timeout.go
func TimeoutHandler(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// cancel to clear resources after finished
			cancel()

			// check if context timeout was reached
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				// 记录操作日志
				logger.Info(ctx, "server timeout", map[string]interface{}{
					"trace_error": ctx.Err(),
				})

				// write response and abort the request
				if !c.IsAborted() {
					c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
						"code":    504,
						"message": http.StatusText(http.StatusGatewayTimeout),
					})
				}
			}
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
