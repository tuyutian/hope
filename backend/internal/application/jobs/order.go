package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"

	"backend/internal/domain/entity/jobs"
	"backend/internal/domain/entity/orders"
	"backend/internal/domain/entity/shopifys"
	jobRepo "backend/internal/domain/repo/jobs"
	orderRepo "backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/products"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/infras/task"
	"backend/internal/providers"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type OrderService struct {
	orderRepo        orderRepo.OrderRepository
	orderInfoRepo    orderRepo.OrderInfoRepository
	orderSummaryRepo orderRepo.OrderSummaryRepository
	jobOrderRepo     jobRepo.OrderRepository
	userRepo         users.UserRepository
	variantRepo      products.VariantRepository
	orderGraphqlRepo shopifyRepo.OrderGraphqlRepository
	shopifyRepo      shopifyRepo.ShopifyRepository
}

func NewOrderService(repos *providers.Repositories) *OrderService {
	return &OrderService{
		orderRepo:        repos.OrderRepo,
		orderInfoRepo:    repos.OrderInfoRep,
		orderSummaryRepo: repos.OrderSummaryRepo,
		jobOrderRepo:     repos.JobOrderRepo,
		userRepo:         repos.UserRepo,
		orderGraphqlRepo: repos.OrderGraphqlRepo,
		shopifyRepo:      repos.ShopifyRepo,
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
	var insuranceAmount float64

	for _, lineItem := range data.Order.LineItems.Edges {
		variantID := o.shopifyRepo.GetIdFromShopifyGraphqlId(lineItem.Node.Variant.ID)
		refundQuantity := refundMap[variantID]
		skuNum++

		if _, exists := existingVariantMap[variantID]; exists {
			// 更新退款数量
			o.orderInfoRepo.UpdateShopifyVariants(ctx, dbOrderId, variantID, &orders.UserOrderInfo{RefundNum: refundQuantity})
		} else {
			// 新增订单详情
			price := o.parsePrice(lineItem.Node.OriginalUnitPriceSet.ShopMoney.Amount)
			isInsurance := 0
			if _, ok := variantIDMap[variantID]; ok {
				isInsurance = 1
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
				IsInsurance:     isInsurance,
				UserOrderId:     dbOrderId,
			})
		}
	}

	userOrder.SkuNum = skuNum
	userOrder.InsuranceAmount = insuranceAmount
	o.orderRepo.UpdateShopifyOrderId(ctx, userOrder)

	// 插入新增的变体
	if len(userOrderInfos) > 0 {
		return o.orderInfoRepo.Create(ctx, userOrderInfos)
	}

	return nil
}

func (o *OrderService) createNewOrder(ctx context.Context, userID int64, data *shopifys.OrderResponse, variantIDMap map[int64]struct{}) error {
	refundMap, refundAmount := o.parseRefundInfo(data)

	total := o.parsePrice(data.Order.TotalPriceSet.ShopMoney.Amount)
	createdAt := utils.PaseTimeToStamp(data.Order.CreatedAt)
	processedAt := utils.PaseTimeToStamp(data.Order.ProcessedAt)

	userOrder := &orders.UserOrder{
		UserID:            userID,
		OrderId:           o.shopifyRepo.GetIdFromShopifyGraphqlId(data.Order.ID),
		OrderName:         data.Order.Name,
		OrderCreatedAt:    createdAt,
		OrderCompletionAt: processedAt,
		FinancialStatus:   data.Order.DisplayFinancialStatus,
		TotalPriceAmount:  total,
		RefundPriceAmount: refundAmount,
		InsuranceAmount:   0,
		Currency:          data.Order.TotalPriceSet.ShopMoney.CurrencyCode,
		SkuNum:            0,
	}

	var userOrderInfos []*orders.UserOrderInfo
	var insuranceAmount float64
	var skuNum int

	for _, lineItem := range data.Order.LineItems.Edges {
		variantID := o.shopifyRepo.GetIdFromShopifyGraphqlId(lineItem.Node.Variant.ID)
		price := o.parsePrice(lineItem.Node.OriginalUnitPriceSet.ShopMoney.Amount)
		refundQuantity := refundMap[variantID]
		isInsurance := 0
		if _, ok := variantIDMap[variantID]; ok {
			isInsurance = 1
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
			IsInsurance:     isInsurance,
		})
		skuNum++
	}

	userOrder.InsuranceAmount = insuranceAmount
	userOrder.SkuNum = skuNum

	dbOrderId, err := o.orderRepo.Create(ctx, userOrder)
	if err != nil {
		return fmt.Errorf("插入主订单失败: %w", err)
	}

	for i := range userOrderInfos {
		userOrderInfos[i].UserOrderId = dbOrderId
	}

	return o.orderInfoRepo.Create(ctx, userOrderInfos)
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
			variantID := o.shopifyRepo.GetIdFromShopifyGraphqlId(item.Node.LineItem.Variant.ID)
			quantity := item.Node.Quantity
			price := o.parsePrice(item.Node.SubtotalSet.ShopMoney.Amount)

			refundMap[variantID] += quantity
			refundAmount += price
		}
	}
	return refundMap, refundAmount
}

func (o *OrderService) parsePrice(amount string) float64 {
	price, _ := strconv.ParseFloat(amount, 64)
	return price
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
	users, err := o.userRepo.GetUsers(ctx, userID, batchSize)
	if err != nil {
		logger.Error(ctx, "order_statistic_queue:获取用户信息失败", err)
		return nil
	}

	if users == nil || len(users) == 0 {
		logger.Info(ctx, "order_statistic_queue:执行完毕")
		return nil
	}

	var lastUserID int64 // 用于存储最后一个处理的用户 userID
	// 处理这一批
	for _, user := range users {
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
			Sales:  statistics.TotalInsurance,
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
		orderStatisticTask, err := task.OrderStatisticsTask(ctx, lastUserID, start, end)
		if err != nil {
			logger.Error(ctx, "order_statistic_queue 构建任务失败:", err.Error())
			return nil
		}
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
