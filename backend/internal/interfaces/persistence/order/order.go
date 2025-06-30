package order

import (
	"context"

	orderEntity "backend/internal/domain/entity/orders"
	orderRepo "backend/internal/domain/repo/orders"

	"xorm.io/xorm"
)

var _ orderRepo.OrderRepository = (*orderRepoImpl)(nil)

type orderRepoImpl struct {
	db *xorm.Engine
}

func NewOrderRepository(db *xorm.Engine) orderRepo.OrderRepository {
	return &orderRepoImpl{db: db}
}

// DelOrder 软删除订单
func (o *orderRepoImpl) DelOrder(ctx context.Context, userID int64, orderId string) error {
	_, err := o.db.Context(ctx).Where("user_id = ? and order_id = ?", userID, orderId).
		Update(&orderEntity.UserOrder{IsDel: 1})
	if err != nil {
		return err
	}
	return nil
}

// List 分页查询订单列表
func (o *orderRepoImpl) List(ctx context.Context, req orderEntity.QueryOrderEntity) ([]*orderEntity.UserOrder, int64, error) {
	var orders []*orderEntity.UserOrder // 直接查询到指针切片，避免二次转换

	// 计算偏移量
	offset := (req.Page - 1) * req.PageSize

	session := o.db.Where("user_id = ? AND is_del = 0", req.UserID)

	// 1. 根据 Type 筛选状态
	switch req.Type {
	case 1:
		session = session.In("financial_status", []string{"PAID", "PARTIALLY_PAID"})
	case 2:
		session = session.In("financial_status", []string{"REFUNDED", "PARTIALLY_REFUNDED"})
		// "All" 不加任何条件
	}

	// 2. Query 模糊搜索（例如搜索订单号）
	if req.Query != "" {
		session = session.Where("order_name LIKE ?", "%"+req.Query+"%")
	}

	// 3. 分页 & 排序
	err := session.
		Desc("id").
		Limit(req.PageSize, offset).
		Find(&orders)

	if err != nil {
		return nil, 0, err
	}
	builder := o.db.Where("user_id = ? AND is_del = 0", req.UserID)

	// 根据 Type 筛选
	switch req.Type {
	case 1:
		builder = builder.In("financial_status", []string{"PAID", "PARTIALLY_PAID"})
	case 2:
		builder = builder.In("financial_status", []string{"REFUNDED", "PARTIALLY_REFUNDED"})
	}

	// 模糊搜索
	if req.Query != "" {
		builder = builder.Where("order_name LIKE ?", "%"+req.Query+"%")
	}

	count, err := builder.Count(&orderEntity.UserOrder{})
	if err != nil {
		return orders, 0, err
	}
	return orders, count, nil
}

// Create 创建订单
func (o *orderRepoImpl) Create(ctx context.Context, order *orderEntity.UserOrder) (int64, error) {
	_, err := o.db.Context(ctx).Insert(order)
	if err != nil {
		return 0, err
	}
	return order.Id, nil
}

// ExistsByOrderID 检查订单是否存在
func (o *orderRepoImpl) ExistsByOrderID(ctx context.Context, orderId int64, userID int64) int64 {
	var userOrder orderEntity.UserOrder
	has, err := o.db.Context(ctx).Cols("id").Where("user_id = ? and order_id = ?", userID, orderId).Get(&userOrder)

	if err != nil || !has {
		return 0
	}

	return userOrder.Id
}

// UpdateShopifyOrderId 更新订单信息
func (o *orderRepoImpl) UpdateShopifyOrderId(ctx context.Context, order *orderEntity.UserOrder) error {
	_, err := o.db.Context(ctx).ID(order.Id).Update(order)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderStatistics 获取订单统计信息
func (o *orderRepoImpl) GetOrderStatistics(ctx context.Context, start, end int64, userID int64) (*orderEntity.OrderStatistics, error) {
	var stats orderEntity.OrderStatistics

	// 在XORM中使用SQL构建统计查询
	has, err := o.db.Context(ctx).SQL("SELECT SUM(refund_price_amount) AS total_refund, SUM(insurance_amount) AS total_insurance, COUNT(*) AS total_orders FROM user_order WHERE order_created_at BETWEEN ? AND ? AND user_id = ?", start, end, userID).Get(&stats)

	if err != nil {
		return nil, err
	}
	if !has {
		return &orderEntity.OrderStatistics{}, nil
	}

	return &stats, nil
}
