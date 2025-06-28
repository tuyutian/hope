package user

import (
	"context"

	"xorm.io/xorm"

	userEntity "backend/internal/domain/entity/users"
	"backend/internal/domain/repo/users"
)

var _ users.AppAuthRepository = (*appAuthRepoImpl)(nil)

type appAuthRepoImpl struct {
	db *xorm.Engine
}

func NewAppAuthRepository(db *xorm.Engine) users.AppAuthRepository {
	return &appAuthRepoImpl{
		db: db,
	}
}

func (a *appAuthRepoImpl) Create(ctx context.Context, appAuth *userEntity.UserAppAuth) (int64, error) {
	_, err := a.db.Context(ctx).Table(userEntity.UserAppAuthTable).Insert(appAuth)
	if err != nil {
		return 0, err
	}

	return appAuth.Id, nil
}

func (a *appAuthRepoImpl) Update(ctx context.Context, appAuth *userEntity.UserAppAuth) error {
	_, err := a.db.Context(ctx).Table(userEntity.UserAppAuthTable).ID(appAuth.Id).Update(appAuth)
	if err != nil {
		return err
	}
	return nil
}

func (a *appAuthRepoImpl) Get(ctx context.Context, id int64, columns ...string) (*userEntity.UserAppAuth, error) {
	var appAuth userEntity.UserAppAuth
	has, err := a.db.Context(ctx).Table(userEntity.UserAppAuthTable).Where("id = ?", id).Cols(columns...).Get(&appAuth)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &appAuth, nil
}
func (a *appAuthRepoImpl) GetByUserAndApp(ctx context.Context, userId int64, appId string, columns ...string) (*userEntity.UserAppAuth, error) {
	var appAuth userEntity.UserAppAuth
	has, err := a.db.Context(ctx).Table(userEntity.UserAppAuthTable).Where("user_id = ? and app_id = ?", userId, appId).Cols(columns...).Get(&appAuth)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &appAuth, nil
}
