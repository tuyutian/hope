package cart

import (
	"context"

	"xorm.io/xorm"

	entity "backend/internal/domain/entity/settings"
	cartRepo "backend/internal/domain/repo/carts"
)

var _ cartRepo.CartSettingRepository = (*cartSettingRepoImpl)(nil)

type cartSettingRepoImpl struct {
	db *xorm.Engine
}

// NewCartSettingRepository 从数据库获取购物车设置资源
func NewCartSettingRepository(engine *xorm.Engine) cartRepo.CartSettingRepository {
	return &cartSettingRepoImpl{db: engine}
}

// First 根据用户ID获取购物车设置
func (s *cartSettingRepoImpl) First(ctx context.Context, userID int64) (*entity.UserCartSetting, error) {
	var cartSetting entity.UserCartSetting
	has, err := s.db.Context(ctx).Where("user_id = ?", userID).Get(&cartSetting)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &cartSetting, nil
}

// Create 创建购物车设置
func (s *cartSettingRepoImpl) Create(ctx context.Context, setting *entity.UserCartSetting) (int64, error) {
	_, err := s.db.Context(ctx).Insert(setting)
	if err != nil {
		return 0, err
	}
	return setting.Id, nil
}

// Update 更新购物车设置
func (s *cartSettingRepoImpl) Update(ctx context.Context, setting *entity.UserCartSetting) error {
	_, err := s.db.Context(ctx).ID(setting.Id).Update(setting)
	if err != nil {
		return err
	}
	return nil
}

// ExistsByShowID 检查是否存在显示购物车的设置
func (s *cartSettingRepoImpl) ExistsByShowID(ctx context.Context, userID int64) int64 {
	exists, err := s.db.Context(ctx).Where("user_id = ? and show_cart = 1", userID).Exist(&entity.UserCartSetting{})

	if err != nil || !exists {
		return 0
	}

	return 1
}

// CloseCart 关闭购物车
func (s *cartSettingRepoImpl) CloseCart(ctx context.Context, userID int64) error {
	zero := 0
	_, err := s.db.Context(ctx).Where("user_id = ?", userID).
		Update(&entity.UserCartSetting{ShowCart: zero}) // ShowCart 设为 0
	if err != nil {
		return err
	}
	return nil
}
