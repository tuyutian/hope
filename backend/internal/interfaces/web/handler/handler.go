package handler

import (
	"backend/internal/application"
	"backend/internal/providers"
)

// Handlers 控制器
type Handlers struct {
	OrderHandler   *OrderHandler
	CommonHandler  *CommonHandler
	UserHandler    *UserHandler
	SettingHandler *SettingHandler
	WebhookHandler *WebHookHandler
	BillingHandler *BillingHandler
}

func InitHandlers(services *application.Services, repos *providers.Repositories) *Handlers {
	orderHandler := NewOrderHandler(services.OrderService)
	commonHandler := NewCommonHandler(repos.AliyunOssRepo)
	userHandler := NewUserHandler(services)
	settingHandler := NewSettingHandler(services)
	webhookHandler := NewWebHookHandler(services)
	billingHandler := NewBillingHandler(services)
	return &Handlers{
		orderHandler,
		commonHandler,
		userHandler,
		settingHandler,
		webhookHandler,
		billingHandler,
	}
}
