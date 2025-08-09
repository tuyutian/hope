package users

import (
	"context"

	"backend/internal/domain/entity"
	billingEntity "backend/internal/domain/entity/billings"
	"backend/internal/domain/repo/billings"
	orderRepo "backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/users"
	"backend/internal/providers"
	"backend/pkg/logger"
)

type BillingService struct {
	commissionBillRepo       billings.CommissionBillRepository
	orderRepo                orderRepo.OrderRepository
	subscriptionRepo         users.UserSubscriptionRepository
	billingPeriodSummaryRepo billings.BillingPeriodSummaryRepository
}

func NewBillingService(repos *providers.Repositories) *BillingService {
	return &BillingService{
		commissionBillRepo:       repos.CommissionBillRepo,
		orderRepo:                repos.OrderRepo,
		billingPeriodSummaryRepo: repos.BillingPeriodSummaryRepo,
		subscriptionRepo:         repos.UserSubscriptionRepo,
	}
}

func (b *BillingService) BillList(ctx context.Context, userID int64, pagination entity.Pagination) *billingEntity.BillingSummaryResponse {
	list, err := b.billingPeriodSummaryRepo.BillingPeriodSummary(ctx, userID, pagination)
	if err != nil {
		logger.Warn(ctx, "query billing period summary error: ", err)
	}
	if list == nil {
		list = make([]*billingEntity.BillingPeriodSummary, 0)
	}
	count, _ := b.billingPeriodSummaryRepo.BillingPeriodCount(ctx, userID)
	return &billingEntity.BillingSummaryResponse{List: list, Total: count}
}

func (b *BillingService) BillDetails(ctx context.Context, userID int64, pagination entity.Pagination) *billingEntity.CommissionListResponse {
	list, err := b.commissionBillRepo.CommissionList(ctx, userID, pagination)
	if err != nil {
		logger.Warn(ctx, "query billing period summary error: ", err)
	}
	if list == nil {
		list = make([]*billingEntity.CommissionBill, 0)
	}
	count, _ := b.commissionBillRepo.CommissionCount(ctx, userID)
	return &billingEntity.CommissionListResponse{List: list, Total: count}
}

func (b *BillingService) CurrentBillDetail(ctx context.Context, userID int64) *billingEntity.CurrentPeriodResponse {
	subscription, err := b.subscriptionRepo.GetActiveSubscription(ctx, userID)
	response := &billingEntity.CurrentPeriodResponse{
		PeriodEnd:   0,
		PeriodStart: 0,
		Amount:      0,
	}
	if err != nil {
		return response
	}

	if subscription == nil {
		return response
	}
	if subscription.CurrentPeriodEnd == 0 {
		return response
	}
	bill, err := b.billingPeriodSummaryRepo.GetByCurrentPeriod(ctx, userID, subscription.CurrentPeriodEnd)
	if bill != nil {
		response.Amount = bill.TotalCommissionAmount
		response.PeriodStart = bill.BillingPeriodStart
		response.PeriodEnd = bill.BillingPeriodEnd
		return response
	} else {
		response.PeriodStart = subscription.CurrentPeriodStart
		response.PeriodEnd = subscription.CurrentPeriodEnd
	}
	return response
}
