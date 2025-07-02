package application

import (
	"backend/internal/application/apps"
	"backend/internal/application/jobs"
	"backend/internal/application/orders"
	"backend/internal/application/products"
	"backend/internal/application/settings"
	"backend/internal/application/users"
	"backend/internal/providers"
)

type Services struct {
	UserService        *users.UserService
	OrderService       *orders.OrderService
	OrderJobService    *jobs.OrderService
	UserJobService     *jobs.UserService
	ProductJobService  *jobs.ProductService
	CartSettingService *settings.CartSettingService
	ProductService     *products.ProductService
	AppService         *apps.AppService
}

func NewServices(repos *providers.Repositories) *Services {
	userService := users.NewUserService(repos)
	orderService := orders.NewOrderService(repos)
	orderJobService := jobs.NewOrderService(repos)
	productJobService := jobs.NewProductService(repos)
	userJobService := jobs.NewUserService(repos)
	cartSettingService := settings.NewCartSettingService(repos)
	productService := products.NewProductService(repos)
	appService := apps.NewAppService(repos)
	return &Services{UserService: userService, OrderService: orderService, OrderJobService: orderJobService, ProductJobService: productJobService, UserJobService: userJobService, CartSettingService: cartSettingService, ProductService: productService, AppService: appService}
}
