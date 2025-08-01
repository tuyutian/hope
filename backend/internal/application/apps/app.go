package apps

import (
	"context"
	"fmt"

	appEntity "backend/internal/domain/entity/apps"
	appRepo "backend/internal/domain/repo/apps"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/logger"
)

type AppService struct {
	appRepo     appRepo.AppRepository
	shopifyRepo shopifyRepo.ShopifyRepository
}

func NewAppService(repos *providers.Repositories) *AppService {
	return &AppService{appRepo: repos.AppRepo, shopifyRepo: repos.ShopifyRepo}
}

func (a *AppService) GetAppConfig(ctx context.Context, appId string) (*appEntity.AppDefinition, error) {
	return a.appRepo.GetByAppId(ctx, appId)
}

func (a *AppService) GetAppID(ctx context.Context) string {
	return ctx.Value(ctxkeys.AppID).(string)
}

func (a *AppService) VerifyWebhook(ctx context.Context, signature string, body []byte) bool {
	appID := ctx.Value(ctxkeys.AppID).(string)
	config, err := a.GetAppConfig(ctx, appID)
	if err != nil {
		logger.Error(ctx, "get app config error: %s", err.Error())
		return false
	}
	appSecret := config.ApiSecret
	// 从配置中获取 webhook secret
	if appSecret == "" {
		fmt.Println("警告: Webhook secret 未配置，跳过签名验证")
		return true // 开发环境可以跳过验证
	}

	return a.shopifyRepo.VerifyWebhook(ctx, appSecret, signature, body)
}
