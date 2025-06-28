package app

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"

	appEntity "backend/internal/domain/entity/apps"
	appRepo "backend/internal/domain/repo/apps"
)

var _ appRepo.AppRepository = (*appRepoImpl)(nil)

type appRepoImpl struct {
	db *xorm.Engine
	rs redis.UniversalClient
}

func NewAppRepository(engine *xorm.Engine, redisClient redis.UniversalClient) appRepo.AppRepository {
	return &appRepoImpl{db: engine, rs: redisClient}
}

func (r *appRepoImpl) GetGoShopifyByAppID(ctx context.Context, appId string) (goshopify.App, error) {
	appConfig, _ := r.GetAppDefinition(ctx, appId)
	if appConfig == nil {
		return goshopify.App{}, nil

	}
	apiKey := appConfig.ApiKey
	appSecret := appConfig.ApiSecret
	redirectUrl := appConfig.CallbackUrl
	scopes := appConfig.Scopes

	return goshopify.App{
		ApiKey:      apiKey,
		ApiSecret:   appSecret,
		RedirectUrl: redirectUrl,
		Scope:       scopes,
	}, nil
}

func (r *appRepoImpl) GetAppDefinition(ctx context.Context, appId string) (*appEntity.AppDefinition, error) {
	var appConfig appEntity.AppDefinition
	err := r.db.Context(ctx).Where("app_id = ?", appId).Find(&appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}
