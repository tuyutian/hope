package orders

import (
	"time"

	"github.com/gin-gonic/gin"

	orderEntity "backend/internal/domain/entity/orders"
	"backend/internal/domain/repo/orders"
	"backend/internal/providers"
	"backend/pkg/logger"
)

type OrderService struct {
	orderRep        orders.OrderRepository
	orderSummaryRep orders.OrderSummaryRepository
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
	Total  float64 `json:"total";o`
}
type OrderSummaryResp struct {
	OrderStatistics      OrderStatistics        `json:"order_statistics"`
	OrderStatisticsTable []OrderStatisticsTable `json:"order_statistics_table"`
}

func (s OrderService) Summary(ctx *gin.Context, userId int64, days int) (interface{}, error) {

	summary, err := s.orderSummaryRep.GetByDays(ctx, userId, days)

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

func (s OrderService) OrderList(ctx *gin.Context, req orderEntity.QueryOrderEntity) (OrderListResp, error) {
	userOrders, count, err := s.orderRep.List(ctx, req)
	return OrderListResp{Orders: userOrders, Count: count}, err
}

func NewOrderService(repos *providers.Repositories) *OrderService {
	return &OrderService{
		orderRep:        repos.OrderRepo,
		orderSummaryRep: repos.OrderSummaryRepo,
	}
}
