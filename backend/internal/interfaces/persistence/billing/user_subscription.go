package billing

import (
	"context"
	"time"

	billingEntity "backend/internal/domain/entity/users"
	"backend/internal/domain/repo/users"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

var _ users.UserSubscriptionRepository = (*userSubscriptionRepoImpl)(nil)

type userSubscriptionRepoImpl struct {
	db *xorm.Engine
}

func NewUserSubscriptionRepository(db *xorm.Engine) users.UserSubscriptionRepository {
	return &userSubscriptionRepoImpl{db: db}
}

// GetActiveSubscription 获取用户活跃订阅
func (u *userSubscriptionRepoImpl) GetActiveSubscription(ctx context.Context, userID int64) (*billingEntity.UserSubscription, error) {
	subscription := &billingEntity.UserSubscription{}
	has, err := u.db.Where("user_id = ? AND subscription_status = ?",
		userID, billingEntity.SubscriptionStatusActive).Get(subscription)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return subscription, nil
}

// UpsertUserSubscription 插入或更新用户订阅
func (u *userSubscriptionRepoImpl) UpsertUserSubscription(ctx context.Context, subscription *billingEntity.UserSubscription) error {
	exists := &billingEntity.UserSubscription{}
	has, err := u.db.Where("user_id = ? AND subscription_id = ?",
		subscription.UserID, subscription.SubscriptionID).Get(exists)
	if err != nil {
		return err
	}

	if has {
		// 更新
		subscription.ID = exists.ID
		subscription.CreateTime = exists.CreateTime
		_, err = u.db.ID(exists.ID).Update(subscription)
		return err
	} else {
		// 插入
		_, err = u.db.Insert(subscription)
		return err
	}
}

// UpdateSubscriptionBalance 更新订阅余额
func (u *userSubscriptionRepoImpl) UpdateSubscriptionBalance(ctx context.Context, id int64, balanceUsed decimal.Decimal) error {
	_, err := u.db.ID(id).Update(&billingEntity.UserSubscription{
		BalanceUsed: balanceUsed,
	})
	return err
}

// GetSubscriptionByLineItemID 根据LineItemID获取订阅
func (u *userSubscriptionRepoImpl) GetSubscriptionByLineItemID(ctx context.Context, lineItemID int64) (*billingEntity.UserSubscription, error) {
	subscription := &billingEntity.UserSubscription{}
	has, err := u.db.Where("subscription_line_item_id = ?", lineItemID).Get(subscription)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return subscription, nil
}

// GetExpiredSubscriptions 获取过期订阅
func (u *userSubscriptionRepoImpl) GetExpiredSubscriptions(ctx context.Context) ([]*billingEntity.UserSubscription, error) {
	subscriptions := make([]*billingEntity.UserSubscription, 0)
	err := u.db.Where("current_period_end < ? AND subscription_status = ?",
		time.Now().Unix(), billingEntity.SubscriptionStatusActive).Find(&subscriptions)
	return subscriptions, err
}
