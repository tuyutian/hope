package application

import (
	"backend/internal/application/orders"
	"backend/internal/application/users"
	"backend/internal/providers"
)

type Services struct {
	UserService  *users.UserService
	OrderService *orders.OrderService
}

func NewServices(repos *providers.Repositories) *Services {
	userService := users.NewUserService(repos)
	orderService := orders.NewOrderService(repos)
	return &Services{userService, orderService}
}
