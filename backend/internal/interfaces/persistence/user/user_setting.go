package user

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity/users"
	userRepo "backend/internal/domain/repo/users"
)

type userSettingRepoImpl struct {
	db *xorm.Engine
}

var _ userRepo.UserSettingRepository = (*userSettingRepoImpl)(nil)

func NewUserSettingRepository(db *xorm.Engine) userRepo.UserSettingRepository {
	return &userSettingRepoImpl{db: db}
}

func (u *userSettingRepoImpl) Get(ctx context.Context, userID int64, name string) (string, error) {
	var setting users.UserSetting
	has, err := u.db.Context(ctx).Table(&setting).Where("user_id = ? and name = ?", userID, name).Get(&setting)
	if err != nil {
		return "", err
	}
	if !has {
		return "", nil
	}
	return setting.Value, nil
}

func (u *userSettingRepoImpl) Set(ctx context.Context, userID int64, name string, value string) error {
	session := u.db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	setting := &users.UserSetting{
		UserId: userID,
		Name:   name,
		Value:  value,
	}

	affected, err := session.Where("user_id = ? AND name = ?", userID, name).Update(setting)
	if err != nil {
		session.Rollback()
		return err
	}

	if affected == 0 {
		if _, err := session.Insert(setting); err != nil {
			session.Rollback()
			return err
		}
	}

	return session.Commit()
}
