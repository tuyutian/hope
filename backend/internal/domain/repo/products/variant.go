package products

import (
	"context"

	"backend/internal/domain/entity/products"
)

type VariantRepository interface {
	// First 根据用户ID获取变体
	First(ctx context.Context, userID int64) (*products.UserVariant, error)
	// FindID 根据用户产品ID查找变体
	FindID(ctx context.Context, userProductId int64) ([]*products.UserVariant, error)
	// CreateVariants 创建产品变体
	CreateVariants(ctx context.Context, variants []*products.UserVariant) error
	// UpdateVariants 更新产品变体
	UpdateVariants(ctx context.Context, id int64, userID int64, variant *products.UserVariant) error
	// GetUploadedVariantIDs 获取已上传的变体ID列表
	GetUploadedVariantIDs(ctx context.Context, userID int64) ([]int64, error)
	// DelShopifyVariant 删除Shopify变体
	DelShopifyVariant(ctx context.Context, userID int64) error
	// GetVariantConfig 获取变体配置
	GetVariantConfig(ctx context.Context, userID int64) (map[string]int64, int64, error)
}
