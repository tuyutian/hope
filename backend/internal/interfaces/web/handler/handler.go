package handler

import (
	"backend/internal/application/orders"
)

// Handlers 控制器
type Handlers struct {
	OrderHandler *OrderHandler
}

func InitHandlers(orderService *orders.OrderService) *Handlers {
	OrderHandler := &OrderHandler{
		orderService: orderService,
	}
	return &Handlers{
		OrderHandler,
	}
}
