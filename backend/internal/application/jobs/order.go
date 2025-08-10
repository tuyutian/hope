package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cartEntity "backend/internal/domain/entity/settings"
	cartSettingRepo "backend/internal/domain/repo/carts"

	"backend/internal/domain/entity/billings"
	"backend/internal/domain/entity/jobs"
	"backend/internal/domain/entity/orders"
	"backend/internal/domain/entity/shopifys"
	billingsRepo "backend/internal/domain/repo/billings"
	jobRepo "backend/internal/domain/repo/jobs"
	orderRepo "backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/products"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/config"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	orderRepo         orderRepo.OrderRepository
	orderInfoRepo     orderRepo.OrderInfoRepository
	orderSummaryRepo  orderRepo.OrderSummaryRepository
	jobOrderRepo      jobRepo.OrderRepository
	userRepo          users.UserRepository
	variantRepo       products.VariantRepository
	orderGraphqlRepo  shopifyRepo.OrderGraphqlRepository
	shopifyRepo       shopifyRepo.ShopifyRepository
	subscriptionRepo  users.UserSubscriptionRepository
	billingPeriodRepo billingsRepo.BillingPeriodSummaryRepository
	cartSettingRepo   cartSettingRepo.CartSettingRepository
}

func NewOrderService(repos *providers.Repositories) *OrderService {
	return &OrderService{
		orderRepo:         repos.OrderRepo,
		orderInfoRepo:     repos.OrderInfoRep,
		orderSummaryRepo:  repos.OrderSummaryRepo,
		jobOrderRepo:      repos.JobOrderRepo,
		userRepo:          repos.UserRepo,
		orderGraphqlRepo:  repos.OrderGraphqlRepo,
		shopifyRepo:       repos.ShopifyRepo,
		subscriptionRepo:  repos.UserSubscriptionRepo,
		billingPeriodRepo: repos.BillingPeriodSummaryRepo,
		variantRepo:       repos.VariantRepo,
		cartSettingRepo:   repos.CartSettingRepo,
	}
}

func (o *OrderService) HandleOrder(ctx context.Context, t *asynq.Task) error {
	var payload jobs.OrderPayload

	defer func() {
		if r := recover(); r != nil {
			_ = o.fail(ctx, payload.JobId, "panic捕获", fmt.Errorf("%v", r))
		}
	}()

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return o.fail(ctx, 0, "payload反序列化失败", err)
	}

	logger.Info(ctx, "order_queue", fmt.Sprintf("开始消费订单任务: %d", payload.JobId))

	return o.processOrderJob(ctx, payload.JobId)
}

func (o *OrderService) processOrderJob(ctx context.Context, jobId int64) error {
	// 查询 Job
	job, err := o.jobOrderRepo.First(ctx, jobId)
	if err != nil || job == nil {
		return o.fail(ctx, jobId, "查询Job失败", err)
	}
	if job.IsSuccess == 1 {
		return o.skip(ctx, job.Id, "任务已完成，跳过")
	}

	if err := o.jobOrderRepo.UpdateJobTime(ctx, job.Id); err != nil {
		return o.fail(ctx, job.Id, "更新Job时间失败", err)
	}

	user, err := o.userRepo.FirstName(ctx, job.Name)
	if err != nil || user == nil || user.IsDel != 0 {
		return o.fail(ctx, job.Id, "查询用户信息失败或卸载", err)
	}

	// 初始化 Shopify client
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	o.orderGraphqlRepo.WithClient(client)
	// 获取订单信息
	data, err := o.orderGraphqlRepo.GetOrderInfo(ctx, job.OrderId)
	if err != nil {
		return o.fail(ctx, job.Id, "拉取Shopify订单信息失败", err)
	}

	validStatuses := map[string]bool{
		"PAID":               true,
		"PARTIALLY_PAID":     true,
		"PARTIALLY_REFUNDED": true,
		"REFUNDED":           true,
	}

	if !validStatuses[data.Order.DisplayFinancialStatus] {
		return o.skip(ctx, job.Id, fmt.Sprintf("订单未支付，跳过: %s", data.Order.DisplayFinancialStatus))
	}

	userID := user.ID

	// 查询已上传的变体
	uploadedVariantIDs, err := o.variantRepo.GetUploadedVariantIDs(ctx, userID)
	if err != nil {
		return o.fail(ctx, job.Id, "查询已上传变体失败", err)
	}

	variantIDMap := o.sliceToMap(uploadedVariantIDs)

	// 处理订单
	dbOrderId := o.orderRepo.ExistsByOrderID(ctx, job.OrderId, userID)
	if dbOrderId > 0 {
		err = o.updateExistingOrder(ctx, dbOrderId, userID, data, variantIDMap)
	} else {
		err = o.createNewOrder(ctx, userID, data, variantIDMap)
	}
	if err != nil {
		return o.fail(ctx, job.Id, "处理订单失败", err)
	}

	return o.ok(ctx, job.Id)
}

// order.go
func (o *OrderService) updateExistingOrder(ctx context.Context, dbOrderId int64, userID int64, data *shopifys.OrderResponse, variantIDMap map[int64]struct{}) error {
	// 获取订单详情里的已有变体
	existingVariantIDs, err := o.orderInfoRepo.GetOrderDetailVariantIDs(ctx, dbOrderId, userID)
	if err != nil {
		return fmt.Errorf("查询订单详情失败: %w", err)
	}
	existingVariantMap := o.sliceToMap(existingVariantIDs)

	// 解析退款
	refundMap, refundAmount := o.parseRefundInfo(data)

	// 更新主订单
	userOrder := &orders.UserOrder{
		Id:                dbOrderId,
		FinancialStatus:   data.Order.DisplayFinancialStatus,
		RefundPriceAmount: refundAmount,
	}

	var userOrderInfos []*orders.UserOrderInfo
	var skuNum int

	var insuranceAmountDecimal decimal.Decimal

	for _, lineItem := range data.Order.LineItems.Edges {
		variantID := utils.GetIdFromShopifyGraphqlId(lineItem.Node.Variant.ID)
		refundQuantity := refundMap[variantID]
		skuNum++

		if _, exists := existingVariantMap[variantID]; exists {
			// 更新退款数量
			o.orderInfoRepo.UpdateShopifyVariants(ctx, dbOrderId, variantID, &orders.UserOrderInfo{RefundNum: refundQuantity})
		} else {
			// 新增订单详情
			price := lineItem.Node.OriginalUnitPriceSet.ShopMoney.Amount
			isProtectify := 0
			if _, ok := variantIDMap[variantID]; ok {
				isProtectify = 1
				insuranceAmountDecimal = insuranceAmountDecimal.Add(price)
			}

			userOrderInfos = append(userOrderInfos, &orders.UserOrderInfo{
				UserID:          userID,
				Sku:             lineItem.Node.Sku,
				VariantId:       variantID,
				VariantTitle:    lineItem.Node.VariantTitle,
				Quantity:        lineItem.Node.Quantity,
				UnitPriceAmount: utils.DecimalToFloat(price), // 转换为 float64 存储
				Currency:        data.Order.TotalPriceSet.ShopMoney.CurrencyCode,
				RefundNum:       refundQuantity,
				IsProtectify:    isProtectify,
				UserOrderId:     dbOrderId,
			})
		}
	}

	userOrder.SkuNum = skuNum
	userOrder.ProtectifyAmount = utils.DecimalToFloat(insuranceAmountDecimal)
	o.orderRepo.UpdateShopifyOrderId(ctx, userOrder)

	// 插入新增的变体
	if len(userOrderInfos) > 0 {
		return o.orderInfoRepo.Create(ctx, userOrderInfos)
	}

	// 更新订单后更新账单相关记录
	if err := o.updateBillingRecords(ctx, userID, userOrder, userOrderInfos); err != nil {
		return fmt.Errorf("更新账单记录失败: %w", err)
	}

	return nil
}

func (o *OrderService) createNewOrder(ctx context.Context, userID int64, data *shopifys.OrderResponse, variantIDMap map[int64]struct{}) error {
	refundMap, refundAmount := o.parseRefundInfo(data)

	total := data.Order.TotalPriceSet.ShopMoney.Amount
	createdAt := utils.PaseTimeToStamp(data.Order.CreatedAt)
	processedAt := utils.PaseTimeToStamp(data.Order.ProcessedAt)

	userOrder := &orders.UserOrder{
		UserID:            userID,
		OrderId:           utils.GetIdFromShopifyGraphqlId(data.Order.ID),
		OrderName:         data.Order.Name,
		OrderCreatedAt:    createdAt,
		OrderCompletionAt: processedAt,
		FinancialStatus:   data.Order.DisplayFinancialStatus,
		TotalPriceAmount:  utils.DecimalToFloat(total),
		RefundPriceAmount: refundAmount,
		ProtectifyAmount:  0,
		Currency:          data.Order.TotalPriceSet.ShopMoney.CurrencyCode,
		SkuNum:            0,
	}

	var userOrderInfos []*orders.UserOrderInfo
	var insuranceAmount float64
	var skuNum int

	for _, lineItem := range data.Order.LineItems.Edges {
		variantID := utils.GetIdFromShopifyGraphqlId(lineItem.Node.Variant.ID)
		price := utils.DecimalToFloat(lineItem.Node.OriginalUnitPriceSet.ShopMoney.Amount)
		refundQuantity := refundMap[variantID]
		isProtectify := 0
		if _, ok := variantIDMap[variantID]; ok {
			isProtectify = 1
			insuranceAmount += price
		}
		userOrderInfos = append(userOrderInfos, &orders.UserOrderInfo{
			UserID:          userID,
			Sku:             lineItem.Node.Sku,
			VariantId:       variantID,
			VariantTitle:    lineItem.Node.VariantTitle,
			Quantity:        lineItem.Node.Quantity,
			UnitPriceAmount: price,
			Currency:        data.Order.TotalPriceSet.ShopMoney.CurrencyCode,
			RefundNum:       refundQuantity,
			IsProtectify:    isProtectify,
		})
		skuNum++
	}

	userOrder.ProtectifyAmount = insuranceAmount
	userOrder.SkuNum = skuNum

	dbOrderId, err := o.orderRepo.Create(ctx, userOrder)
	if err != nil {
		return fmt.Errorf("插入主订单失败: %w", err)
	}

	for i := range userOrderInfos {
		userOrderInfos[i].UserOrderId = dbOrderId
	}

	err = o.orderInfoRepo.Create(ctx, userOrderInfos)

	// 创建订单后更新账单相关记录
	if err := o.updateBillingRecords(ctx, userID, userOrder, userOrderInfos); err != nil {
		return fmt.Errorf("更新账单记录失败: %w", err)
	}

	return err
}

func (o *OrderService) sliceToMap(slice []int64) map[int64]struct{} {
	m := make(map[int64]struct{}, len(slice))
	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}

func (o *OrderService) parseRefundInfo(data *shopifys.OrderResponse) (map[int64]int, float64) {
	refundMap := make(map[int64]int)
	var refundAmount float64
	for _, refund := range data.Order.Refunds {
		for _, item := range refund.RefundLineItems.Edges {
			variantID := utils.GetIdFromShopifyGraphqlId(item.Node.LineItem.Variant.ID)
			quantity := item.Node.Quantity
			price := utils.DecimalToFloat(item.Node.SubtotalSet.ShopMoney.Amount)

			refundMap[variantID] += quantity
			refundAmount += price
		}
	}
	return refundMap, refundAmount
}

func (o *OrderService) HandleOrderStatistics(ctx context.Context, t *asynq.Task) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(ctx, "order_statistic_queue", fmt.Sprintf("panic捕获: %v ", r))
		}
	}()

	var payload jobs.OrderStatisticPayload

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		logger.Error(ctx, "order_statistic_queue:payload 反序列化失败", err)
		return nil
	}

	userID := payload.UserID
	start := payload.Start
	end := payload.End

	batchSize := 300
	// 获取全部userID 一次性300
	userList, err := o.userRepo.GetUsers(ctx, userID, batchSize)
	if err != nil {
		logger.Error(ctx, "order_statistic_queue:获取用户信息失败", err)
		return nil
	}

	if userList == nil || len(userList) == 0 {
		logger.Info(ctx, "order_statistic_queue:执行完毕")
		return nil
	}

	var lastUserID int64 // 用于存储最后一个处理的用户 userID
	// 处理这一批
	for _, user := range userList {
		lastUserID = user.ID
		// 获取订单数据 订单数据 日退款数据
		statistics, err := o.orderRepo.GetOrderStatistics(ctx, start, end, user.ID)

		if err != nil {
			logger.Error(ctx, "order_statistic_queue:查询统计失败", err)
			continue
		}

		orderStatisticId, err := o.orderSummaryRepo.ExistOrder(ctx, user.ID, start)
		if err != nil {
			logger.Error(ctx, "order_statistic_queue:查询每日记录失败", err)
			continue
		}
		orderSummary := orders.OrderSummary{
			Today:  start,
			Orders: statistics.TotalOrders,
			Refund: statistics.TotalRefund,
			Sales:  statistics.TotalProtectify,
		}

		if orderStatisticId > 0 {
			orderSummary.Id = orderStatisticId
			o.orderSummaryRepo.UpsertOrderStatistics(ctx, orderSummary)
		} else {
			orderSummary.UserID = user.ID
			o.orderSummaryRepo.CrateOrderStatistics(ctx, orderSummary)
		}
	}

	// 递归处理下一批用户
	if lastUserID != 0 {
		payload := jobs.OrderStatisticPayload{UserID: userID, Start: start, End: end}
		data, err := json.Marshal(payload)
		if err != nil {
			logger.Error(ctx, "order_statistic_queue 构建任务失败:", err.Error())
			return nil
		}
		orderStatisticTask := asynq.NewTask(config.SendOrderStatistics, data)
		err = o.HandleOrderStatistics(ctx, orderStatisticTask)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (o *OrderService) ok(ctx context.Context, jobId int64) error {
	logger.Info(ctx, "order_queue", fmt.Sprintf("JobId: %d => 成功完成", jobId))
	_ = o.jobOrderRepo.UpdateStatus(ctx, jobId, 1) // 1 表示成功
	return nil
}

func (o *OrderService) fail(ctx context.Context, jobId int64, msg string, err error) error {
	logger.Error(ctx, fmt.Sprintf("order_queue: JobId: %d => %s: %v", jobId, msg, err))
	_ = o.jobOrderRepo.UpdateStatus(ctx, jobId, 3) // 3 表示失败
	return nil
}

func (o *OrderService) skip(ctx context.Context, jobId int64, reason string) error {
	logger.Info(ctx, fmt.Sprintf("order_queue: JobId: %d => 跳过: %s", jobId, reason))
	_ = o.jobOrderRepo.UpdateStatus(ctx, jobId, 2) // 2 表示跳过
	return nil
}

// calculateCommission 根据用户设置计算佣金
func (o *OrderService) calculateCommission(ctx context.Context, userID int64, protectifyAmount float64, orderTotalAmount float64) (float64, float64, error) {
	// 获取用户的购物车设置
	cartSetting, err := o.cartSettingRepo.First(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("获取用户购物车设置失败: %w", err)
	}
	totalAmount := decimal.NewFromFloat(orderTotalAmount)
	var commissionAmount decimal.Decimal
	var commissionRate decimal.Decimal
	// 解析 PricingSelect
	var prices []cartEntity.PriceSelectReq
	if err := json.Unmarshal([]byte(cartSetting.PricingSelect), &prices); err != nil {
		logger.Error(ctx, "get-cart 解析 PricingSelect 失败", "Err:", err.Error())
		return 0, 0, fmt.Errorf("解析 PricingSelect 失败: %w", err)
	}

	// 解析 TiersSelect
	var tiers []cartEntity.TierSelectReq
	if err := json.Unmarshal([]byte(cartSetting.TiersSelect), &tiers); err != nil {
		logger.Error(ctx, "get-cart 解析 TiersSelect 失败", "Err:", err.Error())
		return 0, 0, fmt.Errorf("解析 TiersSelect 失败: %w", err)
	}
	protectifyPrice := decimal.NewFromFloat(protectifyAmount)
	if cartSetting.PricingRule == 0 {
		if cartSetting.PricingType == 0 {
			commissionRate = decimal.NewFromFloat(cartSetting.AllPriceSet).Div(protectifyPrice)
		}
	} else {
		// 按金额计算
		if cartSetting.PricingType == 0 {
			for _, priceRange := range prices {
				minRange := utils.ParseMoneyDecimal(priceRange.Min)
				maxRange := utils.ParseMoneyDecimal(priceRange.Max)
				if totalAmount.GreaterThanOrEqual(minRange) && (totalAmount.LessThanOrEqual(maxRange) || maxRange.Equal(decimal.Zero)) {
					commissionAmount := utils.ParseMoneyDecimal(priceRange.Price)
					// 计算实际费率
					if protectifyAmount > 0 {
						commissionRate = commissionAmount.Div(protectifyPrice)
					}
					break
				}
			}
		} else { // 按比例计算
			for _, tier := range tiers {
				minRange := utils.ParseMoneyDecimal(tier.Min)
				maxRange := utils.ParseMoneyDecimal(tier.Max)
				if totalAmount.GreaterThanOrEqual(minRange) && (totalAmount.LessThanOrEqual(maxRange) || maxRange.Equal(decimal.Zero)) {
					percentage := utils.ParseMoneyDecimal(tier.Percentage)
					commissionRate = percentage.Div(decimal.NewFromInt(100)) // 转换百分比为小数
					commissionAmount = protectifyPrice.Mul(commissionRate)
					break
				}
			}
		}
	}

	// 如果没有找到匹配的价格区间或比例区间
	if commissionAmount == decimal.Zero && commissionRate == decimal.Zero {
		return 0, 0, fmt.Errorf("未找到匹配的价格或比例区间")
	}

	return utils.DecimalToFloat(commissionAmount), utils.DecimalToFloat(commissionRate), nil
}

// updateBillingRecords 更新订单相关的账单记录和周期汇总
func (o *OrderService) updateBillingRecords(ctx context.Context, userID int64, order *orders.UserOrder, orderInfo []*orders.UserOrderInfo) error {
	// 1. 获取用户当前的订阅信息
	subscription, err := o.subscriptionRepo.GetActiveSubscription(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户订阅信息失败: %w", err)
	}

	// 2. 计算佣金金额
	commissionAmount, commissionRate, err := o.calculateCommission(ctx, userID, order.ProtectifyAmount, order.TotalPriceAmount)
	if err != nil {
		return fmt.Errorf("计算佣金失败: %w", err)
	}

	// 3. 计算账单周期
	now := time.Now().Unix()
	billCycle := time.Unix(now, 0).Format("2006-01")
	businessMonth := billCycle

	// 4. 创建commission_bill记录
	_ = &billings.CommissionBill{
		UserId:                userID,
		UserOrderId:           order.Id,
		OrderName:             order.OrderName,
		BillCycle:             billCycle,
		CommissionAmount:      commissionAmount,
		CommissionRate:        commissionRate,
		Currency:              order.Currency,
		SubscriptionId:        subscription.ID,
		OrderProtectifyAmount: order.ProtectifyAmount,
		OrderTotalAmount:      order.TotalPriceAmount,
		ChargeStatus:          0, // 待提交
		CreateTime:            now,
		UpdateTime:            now,
	}

	// 5. 更新billing_period_summary
	periodEnd := time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Unix() - 1
	summary, err := o.billingPeriodRepo.GetByCurrentPeriod(ctx, userID, periodEnd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 创建新的周期汇总
			summary = &billings.BillingPeriodSummary{
				UserId:                userID,
				BillingPeriodEnd:      periodEnd,
				BillingPeriodStart:    time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC).Unix(),
				BillCycle:             billCycle,
				BusinessMonth:         businessMonth,
				Currency:              order.Currency,
				OrderCount:            1,
				BillCount:             1,
				TotalCommissionAmount: commissionAmount,
				PendingAmount:         commissionAmount,
				TotalProtectifyAmount: order.ProtectifyAmount,
				TotalOrderAmount:      order.TotalPriceAmount,
				TotalRefundAmount:     order.RefundPriceAmount,
				Version:               1,
			}
			_, err = o.billingPeriodRepo.CreateBillingPeriodSummary(ctx, summary)
		} else {
			return fmt.Errorf("查询账期汇总失败: %w", err)
		}
	} else {
		// 更新现有周期汇总
		summary.OrderCount += 1
		summary.BillCount += 1
		summary.TotalCommissionAmount += commissionAmount
		summary.PendingAmount += commissionAmount
		summary.TotalProtectifyAmount += order.ProtectifyAmount
		summary.TotalOrderAmount += order.TotalPriceAmount
		summary.TotalRefundAmount += order.RefundPriceAmount
		summary.Version += 1

		err = o.billingPeriodRepo.UpdateBillingPeriodSummary(ctx, summary)
	}
	if err != nil {
		return fmt.Errorf("更新账期汇总失败: %w", err)
	}

	return nil
}
