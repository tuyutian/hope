package apps

import (
	"context"

	appEntity "backend/internal/domain/entity/apps"
)

type AppRepository interface {
	GetByAppId(ctx context.Context, appId string) (*appEntity.AppDefinition, error)
}
