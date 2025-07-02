package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/domain/repo"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
)

type CommonHandler struct {
	response.BaseHandler
	ossRepo repo.AliyunOSSRepository
}

func NewCommonHandler(ossRepo repo.AliyunOSSRepository) *CommonHandler {
	return &CommonHandler{ossRepo: ossRepo}
}

func (c *CommonHandler) Upload(ctx *gin.Context) {
	// 获取前端传递的图片
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	// 拼接uuid的图片名称
	uuid := uuid.New()
	imageName := uuid.String() + file.Filename
	imagePath, err := c.ossRepo.UploadFile(ctx, imageName, file)
	if err != nil {
		c.Error(ctx, code.UploadFailed, message.ErrorBadRequest.Error(), nil)
		return
	}

	c.Success(ctx, "", map[string]interface{}{"imagePath": imagePath})
}
