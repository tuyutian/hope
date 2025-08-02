package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/application/apps"
	appEntity "backend/internal/domain/entity/apps"
	"backend/pkg/ctxkeys"
)

type AppMiddleware struct {
	appService *apps.AppService
}

func NewAppMiddleware(appService *apps.AppService) *AppMiddleware {
	return &AppMiddleware{
		appService: appService,
	}
}

func (m *AppMiddleware) AppMust() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		appId := c.Param("appId")
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
		ctx = context.WithValue(ctx, ctxkeys.AppData, appData)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
