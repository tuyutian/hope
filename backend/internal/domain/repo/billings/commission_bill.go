package billings

import (
	"context"

	"github.com/shopspring/decimal"

	"backend/internal/domain/entity"
	billingEntity "backend/internal/domain/entity/billings"
)

type CommissionBillRepository interface {
	CreateCommission(ctx context.Context, userID int64, orderID int64, amount decimal.Decimal) (int64, error)
	CommissionList(ctx context.Context, userID int64, pagination entity.Pagination) ([]*billingEntity.CommissionBill, error)
	CommissionCount(ctx context.Context, userID int64) (int64, error)
}
