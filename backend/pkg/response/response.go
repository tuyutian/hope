package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	returnCode "backend/pkg/response/code"
)

// EmptyArray 用作空[]返回
type EmptyArray []struct{}

// EmptyObject 空对象{}格式返回
type EmptyObject struct{}

// BaseHandler 基础控制器
type BaseHandler struct{}

func (ctrl *BaseHandler) ajaxReturn(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

// Success returns code,data,message if ctrl response success.
func (ctrl *BaseHandler) Success(ctx *gin.Context, message string, data interface{}) {
	if len([]rune(message)) == 0 {
		message = "ok"
	}

	ctrl.ajaxReturn(ctx, returnCode.Success, message, data)
}

// Error returns code,data,message if ctrl response error.
func (ctrl *BaseHandler) Error(ctx *gin.Context, code int, message string, data interface{}) {
	if code <= 0 {
		code = returnCode.BadRequest
	}
	fmt.Println(ctx.Request.URL.Path, code, message, data)
	ctrl.ajaxReturn(ctx, code, message, data)
}

// Fail returns code,data,message if ctrl response error.
func (ctrl *BaseHandler) Fail(ctx *gin.Context, code int, message string, data interface{}) {
	if code <= 0 {
		code = returnCode.BadRequest
	}

	// 根据业务错误码映射到合适的HTTP状态码
	httpStatus := returnCode.StatusBadRequest
	switch code {
	case returnCode.NotFound, returnCode.AppNotfound, returnCode.UserNotfound:
		httpStatus = returnCode.StatusNotFound
	case returnCode.Unauthorized:
		httpStatus = returnCode.StatusUnauthorized
	case returnCode.ServerOperationFailed, returnCode.ThirdPartInitFailed:
		httpStatus = returnCode.StatusInternalServerError
	default:
		httpStatus = returnCode.StatusBadRequest
	}

	ctx.JSON(httpStatus, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}
