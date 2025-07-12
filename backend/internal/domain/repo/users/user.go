package users

import (
	"context"
	"time"

	"backend/internal/domain/entity/users"
)

type UserRepository interface {
	// FirstName 根据店铺名称查找用户
	FirstName(ctx context.Context, shop string) (*users.User, error)
	// GetUserIDByShop 根据店铺名称获取用户ID
	GetUserIDByShop(ctx context.Context, appId string, shop string) (int64, error)
	// Get 根据 id 获取 users.User
	Get(ctx context.Context, id int64, columns ...string) (*users.User, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *users.User) (int64, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *users.User) error
	// UpdateIsDel 更新用户卸载状态
	UpdateIsDel(ctx context.Context, userID int64) error
	// UpdateIsClose 更新用户关店状态
	UpdateIsClose(ctx context.Context, userID int64, planDisplayName string) error
	// UpdateStep 更新用户引导步骤
	UpdateStep(ctx context.Context, userID int64, steps string) error
	// SetToken 设置用户令牌和密码
	SetToken(ctx context.Context, userID int64, token string, pwd string) error
	// FirstEmail 根据邮箱查找用户
	FirstEmail(ctx context.Context, email string) (*users.User, error)
	// UpdatePublishCollection 更新用户发布集合信息
	UpdatePublishCollection(ctx context.Context, userID int64, publishId string, collection string) error
	// BatchUid 批量获取用户ID
	BatchUid(ctx context.Context, userID int64, batchSize int) ([]*users.User, error)
	// GetByShop 通过域名获取店铺
	GetByShop(ctx context.Context, appId string, shop string) (*users.User, error)
	GetActiveUserByShop(ctx context.Context, appId string, shop string) (*users.User, error)
	GetActiveUser(ctx context.Context, id int64, columns ...string) (*users.User, error)
	GetUsers(ctx context.Context, cursorId int64, size int) ([]*users.User, error)
}

type UserCacheRepository interface {
	// Set 将 users.User 写入缓存
	Set(ctx context.Context, id int64, user *users.User, ttl time.Duration) error
	// Get 从缓存中根据 id 获取 users.User
	Get(ctx context.Context, id int64) (*users.User, error)
	// SetByShop Set 将 users.User 写入缓存
	SetByShop(ctx context.Context, appId string, shop string, user *users.User, ttl time.Duration) error
	// GetByShop Get 从缓存中根据 id 获取 users.User
	GetByShop(ctx context.Context, appId string, shop string) (*users.User, error)
}

type UserSettingRepository interface {
	// Get 获取设置
	Get(ctx context.Context, userID int64, name string) (string, error)
	// Set 设置配置
	Set(ctx context.Context, userID int64, name string, value string) error
}
