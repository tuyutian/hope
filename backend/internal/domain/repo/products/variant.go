package products

import (
	"context"

	"backend/internal/domain/entity/products"
)

type VariantRepository interface {
	// First 根据用户ID获取变体
	First(ctx context.Context, uid int) (*products.UserVariant, error)
	// FindID 根据用户产品ID查找变体
	FindID(ctx context.Context, userProductId int) ([]*products.UserVariant, error)
	// CreateVariants 创建产品变体
	CreateVariants(ctx context.Context, variants []*products.UserVariant) error
	// UpdateVariants 更新产品变体
	UpdateVariants(ctx context.Context, id int, uid int, variant *products.UserVariant) error
	// GetUploadedVariantIDs 获取已上传的变体ID列表
	GetUploadedVariantIDs(ctx context.Context, uid int) ([]string, error)
	// DelShopifyVariant 删除Shopify变体
	DelShopifyVariant(ctx context.Context, uid int) error
	// GetVariantConfig 获取变体配置
	GetVariantConfig(ctx context.Context, uid int) (map[string]string, string, error)
}
