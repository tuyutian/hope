package users

import (
	"context"

	"backend/internal/domain/entity/users"
)

type AppAuthRepository interface {
	Get(ctx context.Context, id int64, columns ...string) (*users.UserAppAuth, error)
	GetByUserAndApp(ctx context.Context, userId int64, appId string, columns ...string) (*users.UserAppAuth, error)
	Create(ctx context.Context, user *users.UserAppAuth) (int64, error)
	Update(ctx context.Context, appAuth *users.UserAppAuth) error
}
