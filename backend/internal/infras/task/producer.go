package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"backend/internal/domain/entity/jobs"
	"backend/internal/infras/config"
	"backend/pkg/logger"
)

func NewProductTask(ctx context.Context, jobId int64, userProductId int64, shopifyProductId int64) (*asynq.Task, error) {
	payload := jobs.ProductPayload{JobId: jobId, UserProductId: userProductId, ShopifyProductId: shopifyProductId}
	logger.Info(ctx, "正在生产上传产品队列")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "NewProductTask生产失败, Error：", err.Error())
		return nil, err
	}

	return asynq.NewTask(config.SendProduct, data), nil
}

func InitWebhookUserTask(ctx context.Context, userID int64) (*asynq.Task, error) {
	// 初始化任务队列
	payload := jobs.InitUserPayload{UserID: userID}
	logger.Info(ctx, "正在初始化用户设置")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "InitWebhookTask生产失败, Error：", err.Error())
		return nil, err
	}
	return asynq.NewTask(config.SendInitUser, data), nil
}

func OrderWebhookTask(ctx context.Context, jobId int64) (*asynq.Task, error) {
	// 初始化任务队列
	payload := jobs.OrderPayload{JobId: jobId}
	logger.Info(ctx, "正在同步订单信息")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "OrderWebhookTask生产失败, Error：", err.Error())
		return nil, err
	}
	return asynq.NewTask(config.SendOrder, data), nil
}

func ProductWebhookUpdateTask(ctx context.Context, userID int64, userProductId int64) (*asynq.Task, error) {
	// 初始化任务队列
	payload := jobs.ShopifyProductPayload{UserID: userID, UserProductId: userProductId}
	logger.Info(ctx, "正在更新用户shopify产品信息")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "ProductWebhookUpdateTask生产失败, Error：", err.Error())
		return nil, err
	}
	return asynq.NewTask(config.SendUpdateProduct, data), nil
}

func OrderStatisticsTask(ctx context.Context, userID int64, start int64, end int64) (*asynq.Task, error) {
	// 初始化任务队列
	payload := jobs.OrderStatisticPayload{UserID: userID, Start: start, End: end}
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "ProductWebhookUpdateTask生产失败, Error：", err.Error())
		return nil, err
	}
	logger.Info(ctx, "正在统计每日订单数据")
	return asynq.NewTask(config.SendOrderStatistics, data), nil
}

func DelProductTask(ctx context.Context, userID int64, productId int64, delType int) (*asynq.Task, error) {
	// 初始化任务队列
	payload := jobs.DelProductPayload{UserID: userID, ProductId: productId, DelType: delType}
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "DelProductTask生产失败, Error：", err.Error())
		return nil, err
	}
	logger.Info(ctx, "正在删除产品相关数据")
	return asynq.NewTask(config.SendDelProduct, data), nil
}
