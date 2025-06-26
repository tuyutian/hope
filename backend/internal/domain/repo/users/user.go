package users

import (
	"context"

	"backend/internal/domain/entity/users"
)

type UserRepository interface {
	// FirstName 根据店铺名称查找用户
	FirstName(ctx context.Context, name string) (*users.User, error)
	// FirstNameByUid 根据店铺名称获取用户ID
	FirstNameByUid(ctx context.Context, name string) (int, error)
	// FirstID 根据ID查找用户
	FirstID(ctx context.Context, id int) (*users.User, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *users.User) (int, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *users.User) error
	// UpdateIsDel 更新用户卸载状态
	UpdateIsDel(ctx context.Context, uid int) error
	// UpdateIsClose 更新用户关店状态
	UpdateIsClose(ctx context.Context, uid int, planDisplayName string) error
	// UpdateStep 更新用户引导步骤
	UpdateStep(ctx context.Context, uid int, steps string) error
	// SetToken 设置用户令牌和密码
	SetToken(ctx context.Context, uid int, token string, pwd string) error
	// FirstEmail 根据邮箱查找用户
	FirstEmail(ctx context.Context, email string) (*users.User, error)
	// UpdatePublishCollection 更新用户发布集合信息
	UpdatePublishCollection(ctx context.Context, uid int, publishId string, collection string) error
	// BatchUid 批量获取用户ID
	BatchUid(ctx context.Context, uid int, batchSize int) ([]*users.User, error)
}
