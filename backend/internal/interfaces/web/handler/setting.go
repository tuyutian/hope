package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/products"
	"backend/internal/application/settings"
	"backend/internal/application/users"
	productEntity "backend/internal/domain/entity/products"
	settingEntity "backend/internal/domain/entity/settings"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

type SettingHandler struct {
	response.BaseHandler
	cartSettingService *settings.CartSettingService
	productService     *products.ProductService
	userService        *users.UserService
}

func NewSettingHandler(services *application.Services) *SettingHandler {
	return &SettingHandler{cartSettingService: services.CartSettingService, productService: services.ProductService, userService: services.UserService}
}

func (s *SettingHandler) GetCart(ctx *gin.Context) {
	claims := s.userService.GetClaims(ctx.Request.Context())
	uid := claims.UserID

	rsp, err := s.cartSettingService.GetCart(ctx, uid)

	if err != nil {
		utils.CallWilding(err.Error())
		return
	}
	s.Success(ctx, "", rsp)
}

func (s *SettingHandler) UpdateCart(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()
	claims := s.userService.GetClaims(ctxWithTrace)

	var settingToggleReq settingEntity.SettingConfigReq
	err := ctx.ShouldBindJSON(&settingToggleReq)

	if err != nil {
		fmt.Println(err.Error())
		s.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}

	uid := claims.UserID

	settingToggleReq.UserID = uid
	productCollection, err := json.Marshal(settingToggleReq.SelectedCollections)
	if err != nil {
		fmt.Println(err.Error())
		s.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), nil)
		return
	}
	// 上传产品操作
	err = s.productService.UploadProduct(ctxWithTrace, &productEntity.ProductReq{
		UserID:      uid,
		ProductType: settingToggleReq.ProductTypeInput,
		Collection:  string(productCollection),
	})

	if err != nil {
		utils.CallWilding(err.Error())
		return
	}

	err = s.cartSettingService.SetCartSetting(ctxWithTrace, settingToggleReq)

	if err != nil {
		utils.CallWilding(err.Error())
		return
	}

	s.Success(ctx, "", nil)
}

func (s *SettingHandler) GetPublicCart(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()

	var publicCartReq struct {
		Shop string `form:"shop"`
	}
	err := ctx.ShouldBindJSON(&publicCartReq)

	if err != nil {
		utils.CallWilding(err.Error())

		return
	}

	rsp, err := s.cartSettingService.GetPublicCart(ctxWithTrace, publicCartReq.Shop)

	if err != nil {
		utils.CallWilding(err.Error())
		return
	}

	s.Success(ctx, "", rsp)
}
