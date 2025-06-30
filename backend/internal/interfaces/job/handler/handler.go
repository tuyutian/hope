package handler

import (
	"backend/internal/application"
)

// Handlers 控制器
type Handlers struct {
	ProductHandler *ProductHandler
	UserHandler    *UserHandler
	OrderHandler   *OrderHandler
}

func InitHanders(services *application.Services) *Handlers {
	return &Handlers{
		&ProductHandler{
			productService: services.ProductJobService,
		},
		&UserHandler{
			userService: services.UserJobService,
		},
		&OrderHandler{
			orderService: services.OrderJobService,
		},
	}
}
