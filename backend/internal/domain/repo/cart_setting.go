package repo

import (
	"context"

	"backend/internal/domain/entity"
)

type CartSettingRepo interface {
	// First 根据用户ID获取购物车设置
	First(ctx context.Context, uid int) (*entity.UserCartSetting, error)
	// Create 创建购物车设置
	Create(ctx context.Context, setting *entity.UserCartSetting) (int, error)
	// Update 更新购物车设置
	Update(ctx context.Context, setting *entity.UserCartSetting) error
	// ExistsByShowID 检查是否存在显示购物车的设置
	ExistsByShowID(ctx context.Context, uid int) int
	// CloseCart 关闭购物车
	CloseCart(ctx context.Context, uid int) error
}
