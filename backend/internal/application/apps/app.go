package apps

import (
	"context"

	appEntity "backend/internal/domain/entity/apps"
	appRepo "backend/internal/domain/repo/apps"
	"backend/internal/providers"
)

type AppService struct {
	appRepo appRepo.AppRepository
}

func NewAppService(repos *providers.Repositories) *AppService {
	return &AppService{appRepo: repos.AppRepo}
}

func (a *AppService) GetAppConfig(ctx context.Context, appId string) (*appEntity.AppDefinition, error) {
	return a.appRepo.GetByAppId(ctx, appId)
}
