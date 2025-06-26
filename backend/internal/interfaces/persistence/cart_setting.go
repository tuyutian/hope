package persistence

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity"
	cartRepo "backend/internal/domain/repo/carts"
)

type CartSettingRepImpl struct {
	db *xorm.Engine
}

func NewCartSettingRepository(db *xorm.Engine) cartRepo.CartSettingRepository {
	return &CartSettingRepImpl{db: db}
}

func (s *CartSettingRepImpl) First(ctx context.Context, uid int) (*entity.UserCartSetting, error) {
	var cartSetting entity.UserCartSetting
	has, err := s.db.Context(ctx).Where("uid = ?", uid).Get(&cartSetting)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &cartSetting, nil
}

func (s *CartSettingRepImpl) Create(ctx context.Context, setting *entity.UserCartSetting) (int, error) {
	_, err := s.db.Context(ctx).Insert(setting)
	if err != nil {
		return 0, err
	}
	return setting.Id, nil
}

func (s *CartSettingRepImpl) Update(ctx context.Context, setting *entity.UserCartSetting) error {
	_, err := s.db.Context(ctx).ID(setting.Id).Update(setting)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartSettingRepImpl) ExistsByShowID(ctx context.Context, uid int) int {
	exists, err := s.db.Context(ctx).Where("uid = ? and show_cart = 1", uid).Exist(&entity.UserCartSetting{})

	if err != nil || !exists {
		return 0
	}

	return 1
}

func (s *CartSettingRepImpl) CloseCart(ctx context.Context, uid int) error {
	zero := 0
	_, err := s.db.Context(ctx).Where("uid = ?", uid).
		Update(&entity.UserCartSetting{ShowCart: &zero}) // ShowCart 设为 0
	if err != nil {
		return err
	}
	return nil
}
