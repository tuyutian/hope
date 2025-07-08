package handler

import (
	"github.com/gin-gonic/gin"

	"backend/internal/application/users"
	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
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

type UserStepRequest struct {
	UserID  int64 `json:"user_id"`
	StepKey int   `form:"step_key" binding:"required,oneof=1 2 3 4"`
}

func (u *UserHandler) SetUserStep(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()

	var userStepReq UserStepRequest

	if err := ctx.ShouldBindQuery(&userStepReq); err != nil {
		u.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	claims := ctx.Value(ctxkeys.BizClaims).(jwt.BizClaims)
	userID := claims.UserID
	userStepReq.UserID = userID

	err := u.userService.UpdateUserStep(ctxWithTrace, userStepReq.StepKey)
	if err != nil {
		u.Error(ctx, code.ServerOperationFailed, err.Error(), "")
		return
	}
	u.Success(ctx, "", nil)
}

func (u *UserHandler) GetStep(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()
	uid := u.userService.GetClaims(ctxWithTrace).UserID
	resp, err := u.userService.GetUserStep(ctxWithTrace, uid)
	if err != nil {
		u.Error(ctx, code.ServerOperationFailed, err.Error(), "")
		return
	}
	u.Success(ctx, "", resp)
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
