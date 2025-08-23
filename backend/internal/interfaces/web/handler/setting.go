package handler

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/files"
	"backend/internal/application/products"
	"backend/internal/application/settings"
	"backend/internal/application/users"
	appEntity "backend/internal/domain/entity/apps"
	settingEntity "backend/internal/domain/entity/settings"
	"backend/pkg/ctxkeys"
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
	ctxWithTrace := ctx.Request.Context()
	claims := s.userService.GetClaims(ctxWithTrace)
	uid := claims.UserID

	logger.Info(ctxWithTrace, "开始获取购物车设置", "user_id", uid)

	rsp, err := s.cartSettingService.GetCart(ctxWithTrace, uid)
	if err != nil {
		logger.Error(ctxWithTrace, "获取购物车设置失败", "user_id", uid, "error", err.Error())
		utils.CallWilding(err.Error())
		s.Error(ctx, code.ServerOperationFailed, "获取购物车设置失败", nil)
		return
	}

	logger.Info(ctxWithTrace, "获取购物车设置成功", "user_id", uid)
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

	err = s.cartSettingService.SetCartSetting(ctxWithTrace, settingToggleReq)

	if err != nil {
		utils.CallWilding(err.Error())
		return
	}
	var selectIconUrl string
	for _, icon := range settingToggleReq.Icons {
		if icon.Selected {
			selectIconUrl = icon.Src
		}
	}
	if selectIconUrl != "" {
		// 上传产品操作
		err = s.productService.UploadProduct(ctxWithTrace, uid, selectIconUrl)
	}

	if err != nil {
		utils.CallWilding(err.Error())
		logger.Error(ctx, "Update product error: ", err.Error())
		return
	}
	s.Success(ctx, "", nil)
}

func (s *SettingHandler) GetPublicCart(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()
	appData := ctxWithTrace.Value(ctxkeys.AppData).(*appEntity.AppData)
	var publicCartReq struct {
		Shop string `form:"shop"`
	}
	err := ctx.ShouldBindJSON(&publicCartReq)

	if err != nil {
		utils.CallWilding(err.Error())

		return
	}

	rsp, err := s.cartSettingService.GetPublicCart(ctxWithTrace, appData.AppID, publicCartReq.Shop)

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

	imageMedia, err := s.fileService.UploadProductImageToShopify(ctx, image, "Protectify product icon")
	fmt.Println("imageUrl:", imageMedia)
	if err != nil {
		s.Error(c, code.BadRequest, message.ErrUploadFailed.Error(), "")
		return
	}
	s.Success(c, "", struct {
		ID  int64  `json:"id"`
		Src string `json:"src"`
	}{
		ID:  utils.GetIdFromShopifyGraphqlId(imageMedia.ID),
		Src: imageMedia.Image.URL,
	})
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
