package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/application/orders"
	orderEntity "backend/internal/domain/entity/orders"
	"backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
)

type OrderHandler struct {
	orderService *orders.OrderService
	response.BaseHandler
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

	uid, transfer := ctx.Value("id").(int64)
	if !transfer {
		h.Error(ctx, code.Unauthorized, message.ErrorBadRequest.Error(), nil)
		return
	}

	resp, err := h.orderService.Summary(ctx, uid, days)
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

	uid, transfer := ctx.Value("id").(int)
	if !transfer {
		h.Error(ctx, code.Unauthorized, message.ErrorBadRequest.Error(), nil)
		return
	}
	orderListParams.UserID = uid

	resp, err := h.orderService.OrderList(orderListParams)
	if err != nil {
		h.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}
	h.Success(ctx, "", resp)
}
