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

// UpdateSubscriptionStatus Upgrade status by chargeID
func (u *userSubscriptionRepoImpl) UpdateSubscriptionStatus(ctx context.Context, chargeID int64, status string) error {
	_, err := u.db.Where("charge_id = ?", chargeID).Update(&billingEntity.UserSubscription{
		SubscriptionStatus: status,
	})
	if err != nil {
		return err
	}
	return nil
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
	// 只检查记录是否存在
	exists, err := u.db.Where("user_id = ? AND charge_id = ?",
		subscription.UserID, subscription.ChargeID).Exist(&billingEntity.UserSubscription{})
	if err != nil {
		return err
	}

	if exists {
		// 更新
		_, err = u.db.Where("user_id = ? AND charge_id = ?",
			subscription.UserID, subscription.ChargeID).Update(subscription)
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
		BalanceUsed: &balanceUsed,
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

// GetSubscriptionByChargeID 根据LineItemID获取订阅
func (u *userSubscriptionRepoImpl) GetSubscriptionByChargeID(ctx context.Context, chargeID int64) (*billingEntity.UserSubscription, error) {
	subscription := &billingEntity.UserSubscription{}
	has, err := u.db.Where("charge_id = ?", chargeID).Get(subscription)
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

// CancelActiveSubscriptionsExcept 取消用户的所有活跃订阅，除了指定的订阅ID
func (u *userSubscriptionRepoImpl) CancelActiveSubscriptionsExcept(ctx context.Context, userID int64, exceptChargeID int64) error {
	_, err := u.db.Where("user_id = ? AND subscription_status = ? AND charge_id != ?",
		userID, billingEntity.SubscriptionStatusActive, exceptChargeID).
		Update(&billingEntity.UserSubscription{
			SubscriptionStatus: billingEntity.SubscriptionStatusCancelled,
			UpdateTime:         time.Now().Unix(),
		})
	return err
}

// SyncUserSubscriptionWithTx 在事务中同步用户订阅
func (u *userSubscriptionRepoImpl) SyncUserSubscriptionWithTx(ctx context.Context, userID int64, newSubscription *billingEntity.UserSubscription) error {
	// 开启事务
	session := u.db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	// 1. 先插入/更新新订阅
	exists := &billingEntity.UserSubscription{}
	has, err := session.Where("user_id = ? AND charge_id = ?",
		newSubscription.UserID, newSubscription.ChargeID).Get(exists)
	if err != nil {
		session.Rollback()
		return err
	}

	if has {
		// 更新
		newSubscription.ID = exists.ID
		newSubscription.CreateTime = exists.CreateTime
		_, err = session.ID(exists.ID).Update(newSubscription)
	} else {
		// 插入
		_, err = session.Insert(newSubscription)
	}

	if err != nil {
		session.Rollback()
		return err
	}

	// 2. 只有当新订阅是活跃状态时，才取消其他订阅
	if newSubscription.SubscriptionStatus == billingEntity.SubscriptionStatusActive {
		_, err = session.Where("user_id = ? AND subscription_status = ? AND charge_id != ?",
			userID, billingEntity.SubscriptionStatusActive, newSubscription.ChargeID).
			Update(&billingEntity.UserSubscription{
				SubscriptionStatus: billingEntity.SubscriptionStatusCancelled,
				UpdateTime:         time.Now().Unix(),
			})
		if err != nil {
			session.Rollback()
			return err
		}
	}

	// 提交事务
	return session.Commit()
}
