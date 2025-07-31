package handler

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/files"
	"backend/internal/application/products"
	"backend/internal/application/settings"
	"backend/internal/application/users"
	productEntity "backend/internal/domain/entity/products"
	settingEntity "backend/internal/domain/entity/settings"
	"backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

type SettingHandler struct {
	response.BaseHandler
	cartSettingService *settings.CartSettingService
	fileService        *files.FileService
	productService     *products.ProductService
	userService        *users.UserService
}

func NewSettingHandler(services *application.Services) *SettingHandler {
	return &SettingHandler{cartSettingService: services.CartSettingService, productService: services.ProductService, userService: services.UserService, fileService: services.FileService}
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
		UserID:     uid,
		Collection: string(productCollection),
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

func (s *SettingHandler) UploadLogo(c *gin.Context) {
	ctx := c.Request.Context()
	image, err := c.FormFile("image")

	if err != nil {
		logger.Error(ctx, "get image error: ", err.Error())
		s.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	// 定义最大允许大小 (5MB)
	const maxFileSize = 5 << 20 // 5 MiB

	// 验证文件大小
	if image.Size > maxFileSize {
		s.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}

	// 验证文件类型（可选）
	contentType := image.Header.Get("Content-Type")
	fmt.Println(contentType)
	if !isAllowedImageType(contentType) {
		s.Error(c, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}

	imageUrl, err := s.fileService.UploadProductImageToShopify(ctx, image, "Protectify product icon")
	fmt.Println("imageUrl:", imageUrl)
	if err != nil {
		s.Error(c, code.BadRequest, message.ErrUploadFailed.Error(), "")
		return
	}
	s.Success(c, "", imageUrl)
}

// 判断是否允许的图片类型
func isAllowedImageType(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/svg",
		"image/svg+xml",
		"image/webp",
	}
	return slices.Contains(imageTypes, contentType)
}
