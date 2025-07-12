package users

import (
	"context"

	"github.com/shopspring/decimal"

	billingEntity "backend/internal/domain/entity/users"
)

type UserSubscriptionRepository interface {
	GetActiveSubscription(ctx context.Context, userID int64) (*billingEntity.UserSubscription, error)
	UpsertUserSubscription(ctx context.Context, subscription *billingEntity.UserSubscription) error
	UpdateSubscriptionBalance(ctx context.Context, id int64, balanceUsed decimal.Decimal) error
	GetSubscriptionByLineItemID(ctx context.Context, lineItemID int64) (*billingEntity.UserSubscription, error)
	GetExpiredSubscriptions(ctx context.Context) ([]*billingEntity.UserSubscription, error)
}
