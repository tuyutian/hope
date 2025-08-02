package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/application/apps"
	appEntity "backend/internal/domain/entity/apps"
	"backend/internal/domain/repo/jwtauth"
	"backend/internal/infras/config"

	"backend/pkg/ctxkeys"
	"backend/pkg/jwt"
	"backend/pkg/logger"
)

type AppMiddleware struct {
	appService *apps.AppService
	jwtRepo    jwtauth.JWTRepository
	conf       config.JWT
}

func NewAppMiddleware(appService *apps.AppService, jwtRepo jwtauth.JWTRepository, conf config.JWT) *AppMiddleware {
	return &AppMiddleware{
		appService: appService,
		jwtRepo:    jwtRepo,
		conf:       conf,
	}
}

func (m *AppMiddleware) AppMust() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		appId := c.Param("appId")
		logger.Warn(ctx, "AppMiddleware AppMust", zap.String("appId", appId))
		if appId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "appId is empty",
			})
			return
		}
		appDefinition, err := m.appService.GetAppConfig(ctx, appId)
		if err != nil || appDefinition == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "app is empty",
			})
			return
		}
		appData := &appEntity.AppData{
			AppID:     appDefinition.AppId,
			AppKey:    appDefinition.ApiKey,
			AppSecret: appDefinition.ApiSecret,
		}
		jwtManager := jwt.New(
			appDefinition.ApiSecret,
			jwt.WithAccessExpiration(m.conf.AccessExpiration),
			jwt.WithRefreshExpiration(m.conf.RefreshExpiration),
		)
		// 更新现有的 JwtRepository 而不是重新创建
		m.jwtRepo.UpdateSecret(appDefinition.ApiSecret, jwtManager)
		logger.Warn(ctx, "appDefinition dump", zap.Any("appDefinition", appDefinition))
		// 添加日志验证 JwtRepo 已更新
		logger.Info(ctx, "JwtRepo updated for app",
			"appId", appId,
			"appSecret", appDefinition.ApiSecret[:8]+"...")

		ctx = context.WithValue(ctx, ctxkeys.AppData, appData)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
