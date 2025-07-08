package billing

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"

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
