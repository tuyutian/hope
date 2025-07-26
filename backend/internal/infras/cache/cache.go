package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	appEntity "backend/internal/domain/entity/apps"
	"backend/internal/domain/repo"
)

var _ repo.CacheRepository = (*cacheRepoImpl)(nil)

type cacheRepoImpl struct {
	redisClient redis.UniversalClient
}

func NewCacheRepository(redisClient redis.UniversalClient) repo.CacheRepository {
	return &cacheRepoImpl{redisClient}
}

func (c *cacheRepoImpl) GetAppDefinition(ctx context.Context, cacheKey string) (*appEntity.AppDefinition, error) {
	appDefinitionCache, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	var appDef appEntity.AppDefinition
	err = json.Unmarshal([]byte(appDefinitionCache), &appDef)
	if err != nil {
		return nil, err
	}

	return &appDef, nil
}
