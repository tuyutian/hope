package product

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity/products"
	productRepo "backend/internal/domain/repo/products"
)

var _ productRepo.VariantRepository = (*variantRepoImpl)(nil)

type variantRepoImpl struct {
	db *xorm.Engine
}

// NewVariantRepository 从数据库获取产品变体资源
func NewVariantRepository(engine *xorm.Engine) productRepo.VariantRepository {
	return &variantRepoImpl{db: engine}
}

func (v *variantRepoImpl) First(ctx context.Context, userID int64) (*products.UserVariant, error) {
	var variant products.UserVariant
	has, err := v.db.Context(ctx).Where("user_id = ?", userID).Get(&variant)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &variant, nil
}

func (v *variantRepoImpl) FindID(ctx context.Context, userProductId int64) ([]*products.UserVariant, error) {
	var variants []*products.UserVariant
	err := v.db.Context(ctx).Where("user_product_id = ?", userProductId).Find(&variants)

	if err != nil {
		return nil, err
	}
	if len(variants) == 0 {
		return nil, nil
	}

	return variants, nil
}

func (v *variantRepoImpl) CreateVariants(ctx context.Context, variants []*products.UserVariant) error {
	_, err := v.db.Context(ctx).Insert(&variants)
	if err != nil {
		return err
	}
	return nil
}

func (v *variantRepoImpl) UpdateVariants(ctx context.Context, id int64, userID int64, variant *products.UserVariant) error {
	_, err := v.db.Context(ctx).Where("id = ? and user_id = ?", id, userID).Update(variant)
	if err != nil {
		return err
	}
	return nil
}

func (v *variantRepoImpl) GetUploadedVariantIDs(ctx context.Context, userID int64) ([]int64, error) {
	var variantIDs []int64

	err := v.db.Context(ctx).
		Table(new(products.UserVariant)).
		Where("user_id = ?", userID).
		Cols("variant_id").
		Find(&variantIDs)
	if err != nil {
		return nil, err
	}

	return variantIDs, nil
}

func (v *variantRepoImpl) DelShopifyVariant(ctx context.Context, userID int64) error {
	// 使用XORM的Update方法更新多个字段
	_, err := v.db.Context(ctx).
		Where("user_id = ?", userID).
		Update(&products.UserVariant{
			ProductId:   0,
			VariantId:   0,
			InventoryId: 0,
		})
	if err != nil {
		return err
	}
	return nil
}

func (v *variantRepoImpl) GetVariantConfig(ctx context.Context, userID int64) (map[string]int64, int64, error) {
	var userVariants []products.UserVariant

	// 查询所有数据
	err := v.db.Context(ctx).
		Where("user_id = ? AND user_product_id != ?", userID, "").
		Cols("sku_name", "product_id", "variant_id").
		Find(&userVariants)
	if err != nil {
		return nil, 0, err
	}

	resultMap := make(map[string]int64)
	var productId int64

	// 构建 sku_name -> variant_id 的映射
	for _, variant := range userVariants {
		if productId == 0 {
			productId = variant.ProductId
		}
		resultMap[variant.SkuName] = variant.VariantId
	}

	return resultMap, productId, nil
}
