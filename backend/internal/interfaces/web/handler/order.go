package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/application/orders"
	orderEntity "backend/internal/domain/entity/orders"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
)

type OrderHandler struct {
	orderService *orders.OrderService
	response.BaseHandler
}

func NewOrderHandler(orderService *orders.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Dashboard(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	dayParam := strings.TrimSpace(ctx.Query("day"))
	if dayParam == "" {
		dayParam = "30"
	}
	days, err := strconv.Atoi(dayParam)
	// 绑定并校验 query 参数
	if err != nil {
		// 绑定失败或校验失败
		logger.Warn(ctx, "dashboard 参数错误！", "Err:", err.Error())
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	userID := reqCtx.Value(ctxkeys.BizClaims).(*jwt.BizClaims).UserID

	resp, err := h.orderService.Summary(reqCtx, userID, days)
	if err != nil {
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	h.Success(ctx, "", resp)
}

// TestDashboard 公开的测试接口，用于验证 API 是否正常工作
func (h *OrderHandler) TestDashboard(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	dayParam := strings.TrimSpace(ctx.Query("days"))
	if dayParam == "" {
		dayParam = "30"
	}
	days, err := strconv.Atoi(dayParam)
	// 绑定并校验 query 参数
	if err != nil {
		// 绑定失败或校验失败
		logger.Warn(ctx, "test dashboard 参数错误！", "Err:", err.Error())
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	// 使用测试用户 ID (2)
	testUserID := int64(2)

	resp, err := h.orderService.Summary(reqCtx, testUserID, days)
	if err != nil {
		logger.Error(ctx, "test dashboard 服务错误", "Err:", err.Error())
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	h.Success(ctx, "测试接口 - 使用测试用户ID", resp)
}

func (h *OrderHandler) OrderList(c *gin.Context) {
	var orderListParams orderEntity.QueryOrderEntity
	ctx := c.Request.Context()
	// 直接绑定，不要先读取原始数据
	if err := c.ShouldBindJSON(&orderListParams); err != nil {
		fmt.Println("绑定错误:", err.Error())
		h.Error(c, code.BadRequest, "参数错误："+err.Error(), nil)
		return
	}

	claims := ctx.Value(ctxkeys.BizClaims).(*jwt.BizClaims)
	orderListParams.UserID = claims.UserID

	resp, err := h.orderService.OrderList(ctx, orderListParams)
	if err != nil {
		fmt.Println("服务错误:", err.Error())
		h.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}
	h.Success(c, "", resp)
}
