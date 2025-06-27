package product

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity/products"
	productRepo "backend/internal/domain/repo/products"
)

var _ productRepo.ProductRepository = (*productRepoImpl)(nil)

type productRepoImpl struct {
	db *xorm.Engine
}

// NewProductRepository 从数据库获取产品资源
func NewProductRepository(engine *xorm.Engine) productRepo.ProductRepository {
	return &productRepoImpl{db: engine}
}

func (p *productRepoImpl) First(ctx context.Context, userID int64) (*products.UserProduct, error) {
	var product products.UserProduct
	has, err := p.db.Context(ctx).Where("user_id = ?", userID).Get(&product)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &product, nil
}

func (p *productRepoImpl) FirstProductByID(ctx context.Context, id int64, userID int64) (*products.UserProduct, error) {
	var product products.UserProduct
	has, err := p.db.Context(ctx).Where("id = ? and user_id = ?", id, userID).Get(&product)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &product, nil
}

func (p *productRepoImpl) FirstProductID(ctx context.Context, userID int64) string {
	var product products.UserProduct
	has, err := p.db.Context(ctx).Cols("product_id").Where("user_id = ?", userID).Get(&product)

	if err != nil || !has {
		return ""
	}
	return product.ProductId
}

func (p *productRepoImpl) CreateProduct(ctx context.Context, product *products.UserProduct) (int64, error) {
	_, err := p.db.Context(ctx).Insert(product)
	if err != nil {
		return 0, err
	}
	return product.Id, nil
}

func (p *productRepoImpl) UpdateProduct(ctx context.Context, id int64, userID int64, product *products.UserProduct) error {
	_, err := p.db.Context(ctx).Where("id = ? and user_id = ?", id, userID).Update(product)
	if err != nil {
		return err
	}
	return nil
}

func (p *productRepoImpl) DelShopifyProduct(ctx context.Context, userID int64) error {
	_, err := p.db.Context(ctx).Where("user_id = ?", userID).
		Update(&products.UserProduct{ProductId: "", Status: 0})
	if err != nil {
		return err
	}
	return nil
}

func (p *productRepoImpl) ExistsByProductID(ctx context.Context, userID int64, productId string) int64 {
	var userProduct products.UserProduct
	has, err := p.db.Context(ctx).Cols("id").Where("user_id = ? and product_id = ?", userID, productId).Get(&userProduct)

	if err != nil || !has {
		return 0
	}

	return userProduct.Id
}
