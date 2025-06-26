package product

import (
	"backend/internal/domain/entity/products"
	productRepo "backend/internal/domain/repo/products"
	"context"
	"xorm.io/xorm"
)

var _ productRepo.VariantRepository = (*variantRepoImpl)(nil)

type variantRepoImpl struct {
	db *xorm.Engine
}

// NewVariantRepository 从数据库获取产品变体资源
func NewVariantRepository(engine *xorm.Engine) productRepo.VariantRepository {
	return &variantRepoImpl{db: engine}
}

func (v *variantRepoImpl) First(ctx context.Context, uid int) (*products.UserVariant, error) {
	var variant products.UserVariant
	has, err := v.db.Context(ctx).Where("uid = ?", uid).Get(&variant)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &variant, nil
}

func (v *variantRepoImpl) FindID(ctx context.Context, userProductId int) ([]*products.UserVariant, error) {
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

func (v *variantRepoImpl) UpdateVariants(ctx context.Context, id int, uid int, variant *products.UserVariant) error {
	_, err := v.db.Context(ctx).Where("id = ? and uid = ?", id, uid).Update(variant)
	if err != nil {
		return err
	}
	return nil
}

func (v *variantRepoImpl) GetUploadedVariantIDs(ctx context.Context, uid int) ([]string, error) {
	var variantIDs []string

	err := v.db.Context(ctx).
		Table(new(products.UserVariant)).
		Where("uid = ?", uid).
		Cols("variant_id").
		Find(&variantIDs)
	if err != nil {
		return nil, err
	}

	return variantIDs, nil
}

func (v *variantRepoImpl) DelShopifyVariant(ctx context.Context, uid int) error {
	// 使用XORM的Update方法更新多个字段
	_, err := v.db.Context(ctx).
		Where("uid = ?", uid).
		Update(&products.UserVariant{
			ProductId:   "",
			VariantId:   "",
			InventoryId: "",
		})
	if err != nil {
		return err
	}
	return nil
}

func (v *variantRepoImpl) GetVariantConfig(ctx context.Context, uid int) (map[string]string, string, error) {
	var userVariants []products.UserVariant

	// 查询所有数据
	err := v.db.Context(ctx).
		Where("uid = ? AND user_product_id != ?", uid, "").
		Cols("sku_name", "product_id", "variant_id").
		Find(&userVariants)
	if err != nil {
		return nil, "", err
	}

	resultMap := make(map[string]string)
	var productId string

	// 构建 sku_name -> variant_id 的映射
	for _, variant := range userVariants {
		if productId == "" {
			productId = variant.ProductId
		}
		resultMap[variant.SkuName] = variant.VariantId
	}

	return resultMap, productId, nil
}
