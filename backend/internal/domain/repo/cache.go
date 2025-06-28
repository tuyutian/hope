package repo

import (
	"context"

	appEntity "backend/internal/domain/entity/apps"
)

type CacheRepository interface {
	GetAppDefinition(ctx context.Context, cacheKey string) (*appEntity.AppDefinition, error)
}
