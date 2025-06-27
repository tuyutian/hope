package user

import (
	"context"
	"time"

	userRepo "backend/internal/domain/repo/users"

	"xorm.io/xorm"

	"backend/internal/domain/entity/users"
)

var _ userRepo.UserRepository = (*userRepoImpl)(nil)

type userRepoImpl struct {
	db *xorm.Engine
}

// NewUserRepository 从数据库获取用户资源
func NewUserRepository(engine *xorm.Engine) userRepo.UserRepository {
	return &userRepoImpl{db: engine}
}

// FirstName 根据店铺名称查找用户
func (u *userRepoImpl) FirstName(ctx context.Context, name string) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("name = ?", name).Get(&user)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

// FirstNameByUid 根据店铺名称获取用户ID
func (u *userRepoImpl) FirstNameByUid(ctx context.Context, name string) (int, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("name = ?", name).Cols("id").Get(&user)

	if err != nil {
		return 0, err
	}
	if !has {
		return 0, nil
	}
	return user.Id, nil
}

// FirstID 根据ID查找用户
func (u *userRepoImpl) FirstID(ctx context.Context, id int) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("id = ?", id).Get(&user)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

// CreateUser 创建用户
func (u *userRepoImpl) CreateUser(ctx context.Context, user *users.User) (int, error) {
	_, err := u.db.Context(ctx).Insert(user)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

// Update 更新用户信息
func (u *userRepoImpl) Update(ctx context.Context, user *users.User) error {
	user.LastLogin = time.Now().Unix()
	_, err := u.db.Context(ctx).ID(user.Id).Update(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateIsDel 更新用户卸载状态
func (u *userRepoImpl) UpdateIsDel(ctx context.Context, userID int) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{IsDel: 2, UnInstallTime: time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

// UpdateIsClose 更新用户关店状态
func (u *userRepoImpl) UpdateIsClose(ctx context.Context, userID int, planDisplayName string) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{PlanDisplayName: planDisplayName, IsDel: 3, UnInstallTime: time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

// UpdateStep 更新用户引导步骤
func (u *userRepoImpl) UpdateStep(ctx context.Context, userID int, steps string) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{Steps: steps})
	if err != nil {
		return err
	}
	return nil
}

// SetToken 设置用户令牌和密码
func (u *userRepoImpl) SetToken(ctx context.Context, userID int, token string, pwd string) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{UserToken: token, Pwd: pwd})
	if err != nil {
		return err
	}
	return nil
}

// FirstEmail 根据邮箱查找用户
func (u *userRepoImpl) FirstEmail(ctx context.Context, email string) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("email = ?", email).Get(&user)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

// UpdatePublishCollection 更新用户发布集合信息
func (u *userRepoImpl) UpdatePublishCollection(ctx context.Context, userID int, publishId string, collection string) error {
	_, err := u.db.Context(ctx).
		Where("id = ?", userID).
		Update(&users.User{PublishId: publishId, Collection: collection})
	if err != nil {
		return err
	}
	return nil
}

// BatchUid 批量获取用户ID
func (u *userRepoImpl) BatchUid(ctx context.Context, userID int, batchSize int) ([]*users.User, error) {
	var users []*users.User
	err := u.db.Context(ctx).Where("id > ? and is_del = 1", userID).Cols("id").Limit(batchSize).Find(&users)

	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users, nil
}
