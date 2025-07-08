package shopifys

import (
	"context"

	shopifyEntity "backend/internal/domain/entity/shopifys"
)

type ProductGraphqlRepository interface {
	BaseGraphqlRepository
	CreateProduct(ctx context.Context, input shopifyEntity.ProductCreateInput) (*shopifyEntity.ProductCreateResponse, error)
	CreateProductWithMedia(ctx context.Context, input shopifyEntity.ProductCreateInput, media shopifyEntity.CreateMediaInput) (*shopifyEntity.ProductCreateResponse, error)
	GetProduct(ctx context.Context, productID string) (*shopifyEntity.ProductResponse, error)
	AddProductImages(ctx context.Context, productID string, images []shopifyEntity.ProductImageInput) error
	DeleteVariant(ctx context.Context, productID int64, variantID int64) error
	CreateVariants(ctx context.Context, productID int64, input []*shopifyEntity.VariantCreateInput) ([]map[string]interface{}, error)
	UpdateProduct(ctx context.Context, productID int64, product shopifyEntity.ProductUpdateInput, media []shopifyEntity.CreateMediaInput) error
	PublishProduct(ctx context.Context, productID int64, publicationID int64) error
	CollectionProduct(ctx context.Context, shopifyProductID, collectionID int64) error
	UpdateVariants(ctx context.Context, productId int64, variants []*shopifyEntity.VariantUpdateInput) error
	GetCollectionList(ctx context.Context) ([]struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}, error)
}
