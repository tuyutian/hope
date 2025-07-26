package apps

import (
	"context"

	"backend/internal/domain/entity/apps"
)

type AppAuthRepository interface {
	Get(ctx context.Context, id int64, columns ...string) (*apps.UserAppAuth, error)
	GetByUserAndApp(ctx context.Context, userId int64, appId string, columns ...string) (*apps.UserAppAuth, error)
	Create(ctx context.Context, user *apps.UserAppAuth) (int64, error)
	Update(ctx context.Context, appAuth *apps.UserAppAuth) error
}
