package handler

import (
	"strconv"

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

	days, err := strconv.Atoi(ctx.Query("days"))
	// 绑定并校验 query 参数
	if err != nil {
		// 绑定失败或校验失败
		logger.Warn(ctx, "dashboard 参数错误！", "Err:", err.Error())
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	userID, transfer := ctx.Value("id").(int64)
	if !transfer {
		h.Error(ctx, code.Unauthorized, message.ErrorBadRequest.Error(), nil)
		return
	}

	resp, err := h.orderService.Summary(ctx, userID, days)
	if err != nil {
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	h.Success(ctx, "", resp)
}

func (h *OrderHandler) OrderList(ctx *gin.Context) {

	var orderListParams orderEntity.QueryOrderEntity

	// 绑定并校验 query 参数
	if err := ctx.ShouldBindQuery(&orderListParams); err != nil {
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}
	claims := ctx.Value(ctxkeys.BizClaims).(*jwt.BizClaims)

	orderListParams.UserID = claims.UserID

	resp, err := h.orderService.OrderList(ctx, orderListParams)
	if err != nil {
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}
	h.Success(ctx, "", resp)
}
