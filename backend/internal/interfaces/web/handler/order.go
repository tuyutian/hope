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

func (h *OrderHandler) OrderList(c *gin.Context) {
	var orderListParams orderEntity.QueryOrderEntity
	ctx := c.Request.Context()
	// 直接绑定，不要先读取原始数据
	if err := c.ShouldBindJSON(&orderListParams); err != nil {
		fmt.Println("绑定错误:", err.Error())
		h.Error(c, code.BadRequest, "参数错误："+err.Error(), nil)
		return
	}

	fmt.Printf("绑定成功: %+v\n", orderListParams)

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
