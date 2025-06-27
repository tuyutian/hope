package order

import (
	"context"

	"xorm.io/xorm"

	"backend/internal/domain/entity/orders"
	orderRepo "backend/internal/domain/repo/orders"
)

var _ orderRepo.OrderInfoRepository = (*infoRepoImpl)(nil)

type infoRepoImpl struct {
	db *xorm.Engine
}

// NewOrderInfoRepository NewInfoRepository 从数据库获取订单详情资源
func NewOrderInfoRepository(engine *xorm.Engine) orderRepo.OrderInfoRepository {
	return &infoRepoImpl{db: engine}
}

func (o *infoRepoImpl) Create(ctx context.Context, orderInfo []*orders.UserOrderInfo) error {
	_, err := o.db.Context(ctx).Insert(orderInfo)
	if err != nil {
		return err
	}
	return nil
}

func (o *infoRepoImpl) UpdateShopifyVariants(ctx context.Context, userOrderId int64, variantId string, orderInfo *orders.UserOrderInfo) error {
	_, err := o.db.Context(ctx).
		Where("user_order_id = ? and variant_id = ?", userOrderId, variantId).
		Update(orderInfo)
	if err != nil {
		return err
	}
	return nil
}

func (o *infoRepoImpl) GetOrderDetailVariantIDs(ctx context.Context, userOrderId int64, userID int64) ([]string, error) {
	var variantIDs []string

	err := o.db.Context(ctx).
		Table(new(orders.UserOrderInfo)).
		Where("user_order_id = ? and user_id = ?", userOrderId, userID).
		Cols("variant_id").
		Find(&variantIDs)
	if err != nil {
		return nil, err
	}

	return variantIDs, nil
}
