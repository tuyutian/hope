package billing

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"

	"backend/internal/domain/entity"
	billingEntity "backend/internal/domain/entity/billings"
	"backend/internal/domain/repo/billings"
)

var _ billings.CommissionBillRepository = (*commissionBillRepoImpl)(nil)

type commissionBillRepoImpl struct {
	db *xorm.Engine
}

func NewCommissionBillRepository(db *xorm.Engine) billings.CommissionBillRepository {
	return &commissionBillRepoImpl{
		db: db,
	}
}

func (c *commissionBillRepoImpl) CreateCommission(ctx context.Context, userID int64, orderID int64, amount decimal.Decimal) (int64, error) {
	// 1. 创建本地账单记录
	bill := &billingEntity.CommissionBill{
		ChargeId:         0,
		UserId:           userID,
		UserOrderId:      orderID,
		CommissionAmount: amount,
		ChargeStatus:     billingEntity.ChargeStatusPending,
		CreateTime:       time.Now().Unix(),
		UpdateTime:       time.Now().Unix(),
	}
	ID, err := c.db.Table(bill.TableName()).Insert(bill)
	// 2. 保存到数据库
	if err != nil {
		return 0, err
	}
	return ID, nil

}

func (c *commissionBillRepoImpl) UpdateCommission(ctx context.Context, id int64, chargeId int64, chargeStatus string) error {
	_, err := c.db.Table(new(billingEntity.CommissionBill)).
		Where("id = ?", id).
		Update(map[string]interface{}{
			"charge_id":     chargeId,
			"charge_status": chargeStatus,
			"update_time":   time.Now().Unix(),
		})
	return err
}

func (c *commissionBillRepoImpl) GetCommission(ctx context.Context, id int64) (*billingEntity.CommissionBill, error) {
	var bill billingEntity.CommissionBill
	has, err := c.db.Table(new(billingEntity.CommissionBill)).
		Where("id = ?", id).
		Get(&bill)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &bill, nil
}

func (c *commissionBillRepoImpl) CommissionList(ctx context.Context, userID int64, pagination entity.Pagination) ([]*billingEntity.CommissionBill, error) {
	var bills []*billingEntity.CommissionBill
	err := c.db.Table(new(billingEntity.CommissionBill)).
		Where("user_id = ?", userID).
		Desc("create_time").
		Limit(pagination.Size, (pagination.Page-1)*pagination.Size).
		Find(&bills)
	return bills, err
}

func (c *commissionBillRepoImpl) CommissionCount(ctx context.Context, userID int64) (int64, error) {
	var bill billingEntity.CommissionBill
	count, err := c.db.Table(&bill).
		Where("user_id = ?", userID).
		Count(bill)
	if err != nil {
		return 0, err
	}
	return count, nil
}
