package billings

import (
	"context"

	"github.com/shopspring/decimal"
)

type CommissionBillRepository interface {
	CreateCommission(ctx context.Context, userID int64, orderID int64, amount decimal.Decimal) (int64, error)
}
