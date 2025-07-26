package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"backend/internal/application"
	"backend/internal/application/users"
	userEntity "backend/internal/domain/entity/users"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
)

type UserHandler struct {
	response.BaseHandler
	userService         *users.UserService
	subscriptionService *users.SubscriptionService
}

func NewUserHandler(services *application.Services) *UserHandler {
	return &UserHandler{userService: services.UserService, subscriptionService: services.SubscriptionService}
}

func (u *UserHandler) SetUserStep(c *gin.Context) {
	ctx := c.Request.Context()

	var userStepReq userEntity.UpdateStep

	if err := c.ShouldBindJSON(&userStepReq); err != nil {
		u.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	err := u.userService.UpdateUserStep(ctx, userStepReq)
	if err != nil {
		u.Error(c, code.ServerOperationFailed, err.Error(), "")
		return
	}
	u.Success(c, "", nil)
}
func (u *UserHandler) UpdateUserSetting(c *gin.Context) {
	ctx := c.Request.Context()

	var userSetting userEntity.UpdateSetting

	if err := c.ShouldBindJSON(&userSetting); err != nil {
		u.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	err := u.userService.UpdateUserSetting(ctx, userSetting)
	if err != nil {
		u.Error(c, code.ServerOperationFailed, err.Error(), "")
		return
	}
	u.Success(c, "", nil)
}

func (u *UserHandler) GetUserConf(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()

	uid := u.userService.GetClaims(ctxWithTrace).UserID

	resp, err := u.userService.GetUserConf(ctxWithTrace, uid)
	if err != nil {
		u.Error(ctx, code.ServerOperationFailed, err.Error(), "")
		return
	}
	u.Success(ctx, "", resp)
}

func (u *UserHandler) GetSessionData(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()
	uid := u.userService.GetClaims(ctxWithTrace).UserID
	resp, err := u.userService.GetSessionData(ctxWithTrace, uid)
	if err != nil {
		u.Error(ctx, code.ServerOperationFailed, err.Error(), "")
		return
	}

	u.Success(ctx, "", resp)
}

func (u *UserHandler) CreateSubscribe(c *gin.Context) {
	ctx := c.Request.Context()
	subscription, confirmUrl, err := u.subscriptionService.CreateUsageSubscription(
		ctx,
		"Insurance tax",
		decimal.NewFromInt(200),
		"USD",
		"every paid order with insurance product will be taxed",
		true,
	)
	if err != nil {
		u.Error(c, code.PaymentRequestFailed, err.Error(), "")
		return
	}
	if subscription == nil {
		u.Error(c, code.PaymentRequestFailed, "subscription is nil", "")
		return
	}
	u.Success(c, "", confirmUrl)
}
