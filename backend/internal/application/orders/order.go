package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/domain/entity/jobs"
	orderEntity "backend/internal/domain/entity/orders"
	jobRepo "backend/internal/domain/repo/jobs"
	"backend/internal/domain/repo/orders"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/providers"
	"backend/pkg/logger"
)

type OrderService struct {
	orderRepo       orders.OrderRepository
	orderSummaryRep orders.OrderSummaryRepository
	jobOrderRepo    jobRepo.OrderRepository
	userRepo        userRepo.UserRepository
	asynqRepo       jobRepo.AsynqRepository
}
type OrderStatisticsTable struct {
	Date   string  `json:"date"`
	Sales  float64 `json:"sales"`
	Refund float64 `json:"refund"`
}

type OrderStatistics struct {
	Orders int     `json:"orders"`
	Sales  float64 `json:"sales"`
	Refund float64 `json:"refund"`
	Total  float64 `json:"total,omitempty"`
}
type OrderSummaryResp struct {
	OrderStatistics      OrderStatistics        `json:"order_statistics"`
	OrderStatisticsTable []OrderStatisticsTable `json:"order_statistics_table"`
}

func NewOrderService(repos *providers.Repositories) *OrderService {
	return &OrderService{
		jobOrderRepo:    repos.JobOrderRepo,
		orderRepo:       repos.OrderRepo,
		orderSummaryRep: repos.OrderSummaryRepo,
		asynqRepo:       repos.AsyncRepo,
		userRepo:        repos.UserRepo,
	}
}

func (o *OrderService) Summary(ctx *gin.Context, userId int64, days int) (interface{}, error) {

	summary, err := o.orderSummaryRep.GetByDays(ctx, userId, days)

	if err != nil {
		logger.Error(ctx, "summary-db异常:"+err.Error())
		return nil, err
	}

	// 如果没有记录，直接返回默认的统计数据
	if len(summary) == 0 {
		return &OrderSummaryResp{}, nil
	}

	// 2. 加载美国时区
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		logger.Error(ctx, "summary-加载时间异常:"+err.Error())
		return nil, err
	}

	var orderSummaryResp OrderSummaryResp

	for _, v := range summary {
		// 将 v.Today (时间戳) 转换为 time.Time
		t := time.Unix(int64(v.Today), 0)

		// 转换到美国时区
		t = t.In(loc)

		// 格式化为美国时区的 Y-m-d 字符串
		dateStr := t.Format("2006-01-02")

		orderSummaryResp.OrderStatisticsTable = append(orderSummaryResp.OrderStatisticsTable, OrderStatisticsTable{
			Date:   dateStr,
			Sales:  v.Sales,
			Refund: v.Refund,
		})
		orderSummaryResp.OrderStatistics.Sales += v.Sales
		orderSummaryResp.OrderStatistics.Refund += v.Refund
		orderSummaryResp.OrderStatistics.Orders += 1
	}
	orderSummaryResp.OrderStatistics.Total = orderSummaryResp.OrderStatistics.Sales - orderSummaryResp.OrderStatistics.Refund

	return &orderSummaryResp, nil
}

type OrderListResp struct {
	Orders []*orderEntity.UserOrder `json:"orders"`
	Count  int64                    `json:"count"`
}

func (o *OrderService) OrderList(ctx *gin.Context, req orderEntity.QueryOrderEntity) (OrderListResp, error) {
	userOrders, count, err := o.orderRepo.List(ctx, req)
	return OrderListResp{Orders: userOrders, Count: count}, err
}

// OrderSync 处理订单同步 WebHook
func (o *OrderService) OrderSync(ctx context.Context, req orderEntity.OrderWebHookReq) error {
	go func(req orderEntity.OrderWebHookReq) {
		// 每个协程自己 new ctx，超时保护，防止协程卡死
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "OrderSync 协程 panic:", fmt.Sprintf("%v", r))
			}
		}()

		row := o.jobOrderRepo.ExistsByOrderID(newCtx, req.OrderId)
		if row != 0 {
			logger.Info(ctx, "OrderSync 已存在订单记录，无需重复插入：", req.OrderId)
			return
		}

		logger.Info(ctx, "OrderSync 订单日志不存在，开始插入：", req.OrderId, req.Shop)

		log, err := o.jobOrderRepo.Create(newCtx, &jobs.JobOrder{
			OrderId: req.OrderId,
			Name:    req.Shop,
		})
		if err != nil {
			logger.Error(ctx, "OrderSync 插入订单日志失败:", err.Error())
			return
		}

		orderTask, err := o.asynqRepo.OrderWebhookTask(ctx, log)
		if err != nil {
			logger.Error(ctx, "OrderSync 推送订单队列失败:", err.Error(), orderTask)
			return
		}
	}(req)

	return nil
}

// OrderDel 处理订单删除 WebHook
func (o *OrderService) OrderDel(ctx context.Context, req orderEntity.OrderWebHookReq) error {
	go func(req orderEntity.OrderWebHookReq) {
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "OrderDel 协程 panic:", fmt.Sprintf("%v", r))
			}
		}()

		uid, err := o.userRepo.GetUserIDByShop(newCtx, req.AppId, req.Shop)
		if err != nil {
			logger.Error(ctx, "OrderDel 获取UID失败:", err.Error())
			return
		}

		if uid > 0 {
			if err := o.orderRepo.DelOrder(newCtx, uid, req.OrderId); err != nil {
				logger.Info(ctx, "OrderDel 删除订单失败：", req.OrderId, err.Error())
				return
			}
			logger.Info(ctx, "OrderDel 成功删除订单：", req.OrderId)
		}
	}(req)

	return nil
}
