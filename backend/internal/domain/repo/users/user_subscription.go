package users

import (
	"context"

	"github.com/shopspring/decimal"

	billingEntity "backend/internal/domain/entity/users"
)

// 在仓储接口中添加事务支持方法
type UserSubscriptionRepository interface {
	GetActiveSubscription(ctx context.Context, userID int64) (*billingEntity.UserSubscription, error)
	UpsertUserSubscription(ctx context.Context, subscription *billingEntity.UserSubscription) error
	UpdateSubscriptionBalance(ctx context.Context, id int64, balanceUsed decimal.Decimal) error
	GetSubscriptionByLineItemID(ctx context.Context, lineItemID int64) (*billingEntity.UserSubscription, error)
	GetExpiredSubscriptions(ctx context.Context) ([]*billingEntity.UserSubscription, error)
	GetSubscriptionByChargeID(ctx context.Context, chargeID int64) (*billingEntity.UserSubscription, error)
	UpdateSubscriptionStatus(ctx context.Context, chargeID int64, status string) error
	CancelActiveSubscriptionsExcept(ctx context.Context, userID int64, exceptChargeID int64) error
	// 添加事务方法
	SyncUserSubscriptionWithTx(ctx context.Context, userID int64, newSubscription *billingEntity.UserSubscription) error
}
