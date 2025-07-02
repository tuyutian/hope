package app

import (
	"context"

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

func (r *appRepoImpl) GetByAppId(ctx context.Context, appId string) (*appEntity.AppDefinition, error) {
	var appConfig appEntity.AppDefinition
	err := r.db.Context(ctx).Where("app_id = ?", appId).Find(&appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}
