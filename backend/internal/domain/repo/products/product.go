package products

import (
	"context"

	"backend/internal/domain/entity/products"
)

type ProductRepository interface {
	// First 根据用户ID获取产品
	First(ctx context.Context, userID int64) (*products.UserProduct, error)
	// FirstProductByID 根据产品ID和用户ID获取产品
	FirstProductByID(ctx context.Context, id int64, userID int64) (*products.UserProduct, error)
	// FirstProductID 根据用户ID获取产品ID
	FirstProductID(ctx context.Context, userID int64) string
	// CreateProduct 创建产品
	CreateProduct(ctx context.Context, product *products.UserProduct) (int64, error)
	// UpdateProduct 更新产品
	UpdateProduct(ctx context.Context, id int64, userID int64, product *products.UserProduct) error
	// DelShopifyProduct 删除Shopify产品
	DelShopifyProduct(ctx context.Context, userID int64) error
	// ExistsByProductID 根据产品ID检查产品是否存在
	ExistsByProductID(ctx context.Context, userID int64, productId string) int64
}
