package shopifys

import (
	"context"

	shopifyEntity "backend/internal/domain/entity/shopifys"
)

type ProductGraphqlRepository interface {
	BaseGraphqlRepository
	CreateProductWithMedia(ctx context.Context, productInput shopifyEntity.ProductCreateInput, mediaInput []shopifyEntity.CreateMediaInput) (*shopifyEntity.ProductCreateResponse, error)
	GetProduct(ctx context.Context, productID int64) (*shopifyEntity.ProductResponse, error)
	DeleteVariant(ctx context.Context, productID int64, variantID int64) error
	CreateVariants(ctx context.Context, productID int64, input []*shopifyEntity.VariantCreateInput) ([]map[string]interface{}, error)
	UpdateProduct(ctx context.Context, product shopifyEntity.ProductUpdateInput, media []shopifyEntity.CreateMediaInput) (*shopifyEntity.MutationProduct, error)
	UpdateProductComprehensive(ctx context.Context, productID int64, product shopifyEntity.ProductUpdateInput) error
	PublishProduct(ctx context.Context, productID int64, publicationID int64) error
	CollectionProduct(ctx context.Context, shopifyProductID, collectionID int64) error
	UpdateVariants(ctx context.Context, productId int64, variants []*shopifyEntity.VariantUpdateInput) error
	GetCollectionList(ctx context.Context) ([]struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}, error)
	FileCreate(ctx context.Context, input shopifyEntity.FileCreateInput) (*[]shopifyEntity.FileCreated, error)
	StagedUploadsCreate(ctx context.Context, input shopifyEntity.StagedUploadInput) (*[]shopifyEntity.StagedTarget, error)
	GetImageMedia(ctx context.Context, id string) (*shopifyEntity.ImageMedia, error)
	FileUpdate(ctx context.Context, input shopifyEntity.FileUpdateInput) (*[]shopifyEntity.FileUpdated, error)
}
