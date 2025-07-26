package carts

import (
	"context"

	entity "backend/internal/domain/entity/settings"
)

type CartSettingRepository interface {
	// First 根据用户ID获取购物车设置
	First(ctx context.Context, userID int64) (*entity.UserCartSetting, error)
	// Create 创建购物车设置
	Create(ctx context.Context, setting *entity.UserCartSetting) (int64, error)
	// Update 更新购物车设置
	Update(ctx context.Context, setting *entity.UserCartSetting) error
	// ExistsByShowID 检查是否存在显示购物车的设置
	ExistsByShowID(ctx context.Context, userID int64) int64
	// CloseCart 关闭购物车
	CloseCart(ctx context.Context, userID int64) error
}
