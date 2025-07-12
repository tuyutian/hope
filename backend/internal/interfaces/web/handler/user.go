package handler

import (
	"github.com/gin-gonic/gin"

	"backend/internal/application/users"
	userEntity "backend/internal/domain/entity/users"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
)

type UserHandler struct {
	response.BaseHandler
	userService *users.UserService
}

func NewUserHandler(userService *users.UserService) *UserHandler {
	return &UserHandler{userService: userService}
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
