package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"

	appEntity "backend/internal/domain/entity/apps"
	appRepo "backend/internal/domain/repo/apps"
	"backend/pkg/logger"
)

var _ appRepo.AppRepository = (*appRepoImpl)(nil)

type appRepoImpl struct {
	db *xorm.Engine
	rs redis.UniversalClient
}

func NewAppRepository(engine *xorm.Engine, redisClient redis.UniversalClient) appRepo.AppRepository {
	return &appRepoImpl{db: engine, rs: redisClient}
}

func (r *appRepoImpl) GetByAppId(ctx context.Context, appId string) (*appEntity.AppDefinition, error) {
	var appConfig = &appEntity.AppDefinition{}
	has, err := r.db.Context(ctx).Where("app_id = ?", appId).Get(appConfig)
	if err != nil {
		fmt.Println(err)
		logger.Error(ctx, "get app config error: %s", err.Error())
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return appConfig, nil
}
