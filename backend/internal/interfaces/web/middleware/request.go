package middleware

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"backend/pkg/ctxkeys"
	"backend/pkg/logger"
	"backend/pkg/response/code"
	"backend/pkg/utils"
)

// RequestWare 中间件
type RequestWare struct{}

// Access 记录访问日志
func (ware *RequestWare) Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// uri := c.Request.RequestURI

		// 性能分析后发现log.Println输出需要分配大量的内存空间,而且每次写入都需要枷锁处理
		// log.Println("request before")
		// log.Println("request uri: ", uri)

		// 如果采用了nginx x-request-id功能，可以获得x-request-id
		logId := c.GetHeader("X-Request-ID")
		if logId == "" {
			logId = utils.RndUuid() // 日志id
		}

		// 设置跟请求相关的ctx信息
		ip := c.ClientIP()
		requestURI := c.Request.RequestURI
		ua := c.GetHeader("User-Agent")
		method := c.Request.Method
		c.Request = utils.SetValueToHTTPCtx(c.Request, ctxkeys.XRequestID, logId)
		c.Request = utils.SetValueToHTTPCtx(c.Request, ctxkeys.ReqClientIP, ip)
		c.Request = utils.SetValueToHTTPCtx(c.Request, ctxkeys.RequestURI, requestURI)
		c.Request = utils.SetValueToHTTPCtx(c.Request, ctxkeys.UserAgent, ua)
		c.Request = utils.SetValueToHTTPCtx(c.Request, ctxkeys.RequestMethod, method)

		ctx := c.Request.Context()
		// 记录请求日志
		logger.Info(ctx, "exec begin",
			"log_id", logId,
			"request_uri", requestURI,
			"request_method", method,
			"request_agent", ua,
			"request_ip", ip,
		)
		appId := c.Param("appId")
		logger.Warn(ctx, "appid is:"+appId)
		ctx = context.WithValue(ctx, ctxkeys.AppID, appId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		code := c.Writer.Status()
		execTime := time.Now().Sub(t).Seconds() // 计算请求耗时
		ctx2 := c.Request.Context()
		logger.Info(ctx2, "exec end",
			"exec_time", fmt.Sprintf("%.4f", execTime),
			"response_code", code,
			"log_id", logId,
			"request_uri", requestURI,
			"request_method", method,
			"request_agent", ua,
			"request_ip", ip,
		)
	}
}

// Recover gin请求处理中遇到异常或panic捕获
func (ware *RequestWare) Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// log.Printf("error:%v", err)
				ctx := c.Request.Context()
				logger.Warn(ctx, "exec panic",
					"trace_error", fmt.Sprintf("%v", err),
					"full_stack", string(utils.CatchStack()),
				)

				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						errMsg := strings.ToLower(se.Error())
						// 记录操作日志
						logger.Warn(ctx, "os syscall error",
							"trace_error", errMsg,
						)

						if strings.Contains(errMsg, "broken pipe") ||
							strings.Contains(errMsg, "reset by peer") ||
							strings.Contains(errMsg, "request headers: small read buffer") ||
							strings.Contains(errMsg, "unexpected EOF") ||
							strings.Contains(errMsg, "i/o timeout") {
							brokenPipe = true
						}
					}
				}

				// 是否是 brokenPipe类型的错误，如果是，这里就不能往写入流中再写入内容
				// 如果是该类型的错误，就不需要返回任何数据给客户端
				// 代码参考gin recovery.go RecoveryWithWriter方法实现
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					// ctx.Error(err.(error)) // nolint: errcheck
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				// 响应状态
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "server inner error",
				})
			}
		}()

		c.Next()
	}
}

// NotFoundHandler router not found.
func (ware *RequestWare) NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"code":    code.NotFound,
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
