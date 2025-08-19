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

// GetUserIDByShop 根据店铺名称获取用户ID
func (u *userRepoImpl) GetUserIDByShop(ctx context.Context, appId string, name string) (int64, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("name = ?", name).Cols("id").Get(&user)

	if err != nil {
		return 0, err
	}
	if !has {
		return 0, nil
	}
	return user.ID, nil
}

// Get 根据 id 获取 users.User
func (u *userRepoImpl) Get(ctx context.Context, id int64, columns ...string) (*users.User, error) {
	user := &users.User{}
	_, err := u.db.Table(user.TableName()).Cols(columns...).
		Where("id = ?", id).
		Limit(1).
		Get(user)
	return user, err
}

// GetActiveUser 根据 id 获取 users.User
func (u *userRepoImpl) GetActiveUser(ctx context.Context, id int64, columns ...string) (*users.User, error) {
	user := &users.User{}
	_, err := u.db.Table(user.TableName()).Cols(columns...).
		Where("id = ? and is_del = 0", id).
		Limit(1).
		Get(user)
	return user, err
}

// CreateUser 创建用户
func (u *userRepoImpl) CreateUser(ctx context.Context, user *users.User) (int64, error) {
	_, err := u.db.Context(ctx).Insert(user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Update 更新用户信息
func (u *userRepoImpl) Update(ctx context.Context, user *users.User) error {
	user.LastLogin = time.Now().Unix()
	_, err := u.db.Context(ctx).ID(user.ID).Update(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateIsDel 更新用户卸载状态
func (u *userRepoImpl) UpdateIsDel(ctx context.Context, userID int64, isDel int8) error {
	user := &users.User{}
	if isDel > 0 {
		user.UninstallTime = time.Now().Unix()
		user.IsDel = isDel
	} else {
		user.InstallTime = time.Now().Unix()
		user.IsDel = isDel
	}
	_, err := u.db.Context(ctx).Where("id = ?", userID).MustCols("is_del").
		Update(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateIsClose 更新用户关店状态
func (u *userRepoImpl) UpdateIsClose(ctx context.Context, userID int64, planDisplayName string) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{PlanDisplayName: planDisplayName, IsDel: 3, UninstallTime: time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

// UpdateStep 更新用户引导步骤
func (u *userRepoImpl) UpdateStep(ctx context.Context, userID int64, steps string) error {
	// TODO 这块实现要放到 user_setting表里去
	/*_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{Steps: steps})
	if err != nil {
		return err
	}*/
	return nil
}

// SetToken 设置用户令牌和密码
func (u *userRepoImpl) SetToken(ctx context.Context, userID int64, token string, pwd string) error {
	_, err := u.db.Context(ctx).Where("id = ?", userID).
		Update(&users.User{AccessToken: token, Password: pwd})
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
func (u *userRepoImpl) UpdatePublishCollection(ctx context.Context, userID int64, publishId string, collection string) error {
	// TODO 这块实现要放到 user_setting表里去
	/*_, err := u.db.Context(ctx).
		Where("id = ?", userID).
		Update(&users.User{PublishId: publishId, Collection: collection})
	if err != nil {
		return err
	}*/
	return nil
}

// BatchUid 批量获取用户ID
func (u *userRepoImpl) BatchUid(ctx context.Context, userID int64, batchSize int) ([]*users.User, error) {
	var usersList []*users.User
	err := u.db.Context(ctx).Where("id > ? and is_del = 1", userID).Cols("id").Limit(batchSize).Find(&usersList)

	if err != nil {
		return nil, err
	}
	if len(usersList) == 0 {
		return nil, nil
	}
	return usersList, nil
}

// GetByShop 获取用户店铺
func (u *userRepoImpl) GetByShop(ctx context.Context, appId string, shop string) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where("shop = ? and app_id = ?", shop, appId).OrderBy("create_time").Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

// GetActiveUserByShop 获取用户正常店铺
func (u *userRepoImpl) GetActiveUserByShop(ctx context.Context, appId string, shop string) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where(" app_id = ? and shop = ? and is_del = 0", appId, shop).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

// GetActiveUserByShopID GetActiveUserByShop 获取用户正常店铺
func (u *userRepoImpl) GetActiveUserByShopID(ctx context.Context, appId string, shopID int64) (*users.User, error) {
	var user users.User
	has, err := u.db.Context(ctx).Where(" app_id = ? and shop_id = ? and is_del = 0", appId, shopID).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (u *userRepoImpl) GetUsers(ctx context.Context, cursorId int64, limit int) ([]*users.User, error) {
	var usersList []*users.User
	err := u.db.Context(ctx).Where("id > ? and is_del = 0", cursorId).Cols("id").Limit(limit).Find(&usersList)
	if err != nil {
		return nil, err
	}
	if len(usersList) == 0 {
		return nil, nil
	}
	return usersList, nil
}
