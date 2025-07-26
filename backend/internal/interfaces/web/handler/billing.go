package handler

import (
	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/users"
	"backend/internal/domain/entity"
	"backend/pkg/response"
)

type BillingHandler struct {
	response.BaseHandler
	subscriptionService *users.SubscriptionService
	userService         *users.UserService
	billingService      *users.BillingService
}

func NewBillingHandler(services *application.Services) *BillingHandler {
	return &BillingHandler{
		subscriptionService: services.SubscriptionService,
		userService:         services.UserService,
		billingService:      services.BillingService,
	}
}

func (b *BillingHandler) BillList(c *gin.Context) {
	ctx := c.Request.Context()
	claims := b.userService.GetClaims(ctx)
	userID := claims.UserID
	var pagination entity.Pagination
	err := c.ShouldBindJSON(&pagination)
	if err != nil {
		pagination = entity.Pagination{
			Page: 1,
			Size: 20,
		}
	}
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Size <= 0 {
		pagination.Size = 20
	}
	data := b.billingService.BillList(ctx, userID, pagination)

	b.Success(c, "", data)
}

func (b *BillingHandler) BillDetails(c *gin.Context) {
	ctx := c.Request.Context()
	claims := b.userService.GetClaims(ctx)
	userID := claims.UserID
	var pagination entity.Pagination
	err := c.ShouldBindJSON(&pagination)
	if err != nil {
		pagination = entity.Pagination{
			Page: 1,
			Size: 20,
		}
	}
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Size <= 0 {
		pagination.Size = 20
	}
	data := b.billingService.BillDetails(ctx, userID, pagination)

	b.Success(c, "", data)
}

func (b *BillingHandler) CurrentPeriod(c *gin.Context) {
	userID := b.userService.GetClaims(c.Request.Context()).UserID

	data := b.billingService.CurrentBillDetail(c.Request.Context(), userID)

	b.Success(c, "", data)
}
