package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"backend/internal/domain/entity/jobs"
	jobRepo "backend/internal/domain/repo/jobs"
	"backend/internal/infras/config"
	"backend/pkg/logger"
)

var _ jobRepo.AsynqRepository = (*asynqRepoImpl)(nil)

type asynqRepoImpl struct {
	client *asynq.Client
}

func NewAsynqRepository(client *asynq.Client) jobRepo.AsynqRepository {
	return &asynqRepoImpl{client: client}
}

func (a *asynqRepoImpl) NewProductTask(ctx context.Context, jobId int64, userProductId int64, shopifyProductId int64) (*asynq.TaskInfo, error) {
	payload := jobs.ProductPayload{JobId: jobId, UserProductId: userProductId, ShopifyProductId: shopifyProductId}
	logger.Info(ctx, "正在生产上传产品队列")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "NewProductTask生产失败, Error：", err.Error())
		return nil, err
	}

	task := asynq.NewTask(config.SendProduct, data)
	return a.sendEnqueue(ctx, task)
}

func (a *asynqRepoImpl) InitUserTask(ctx context.Context, userID int64) (*asynq.TaskInfo, error) {
	// 初始化任务队列
	payload := jobs.InitUserPayload{UserID: userID}
	logger.Info(ctx, "正在初始化用户设置")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "InitUserTask生产失败, Error：", err.Error())
		return nil, err
	}
	task := asynq.NewTask(config.SendInitUser, data)
	return a.sendEnqueue(ctx, task)
}

func (a *asynqRepoImpl) OrderWebhookTask(ctx context.Context, jobId int64) (*asynq.TaskInfo, error) {
	// 初始化任务队列
	payload := jobs.OrderPayload{JobId: jobId}
	logger.Info(ctx, "正在同步订单信息")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "OrderWebhookTask生产失败, Error：", err.Error())
		return nil, err
	}
	task := asynq.NewTask(config.SendOrder, data)
	return a.sendEnqueue(ctx, task)
}

func (a *asynqRepoImpl) ProductWebhookUpdateTask(ctx context.Context, userID int64, userProductId int64) (*asynq.TaskInfo, error) {
	// 初始化任务队列
	payload := jobs.ShopifyProductPayload{UserID: userID, UserProductId: userProductId}
	logger.Info(ctx, "正在更新用户shopify产品信息")
	// 使用标准库 json.Marshal 进行序列化
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Info(ctx, "ProductWebhookUpdateTask生产失败, Error：", err.Error())
		return nil, err
	}
	task := asynq.NewTask(config.SendUpdateProduct, data)
	return a.sendEnqueue(ctx, task)

}

func (a *asynqRepoImpl) OrderStatisticsTask(ctx context.Context, userID int64, start int64, end int64) (*asynq.TaskInfo, error) {
	// 初始化任务队列
	payload := jobs.OrderStatisticPayload{UserID: userID, Start: start, End: end}
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "order_statistic_queue 构建任务失败:", err.Error())
		return nil, err
	}
	task := asynq.NewTask(config.SendOrderStatistics, data)
	return a.sendEnqueue(ctx, task)
}

func (a *asynqRepoImpl) DelProductTask(ctx context.Context, userID int64, productId int64, delType int) (*asynq.TaskInfo, error) {
	// 初始化任务队列
	payload := jobs.DelProductPayload{UserID: userID, ProductId: productId, DelType: delType}
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "DelProductTask生产失败, Error：", err.Error())
		return nil, err
	}
	logger.Info(ctx, "正在删除产品相关数据")
	task := asynq.NewTask(config.SendDelProduct, data)
	return a.sendEnqueue(ctx, task)
}

func (a *asynqRepoImpl) sendEnqueue(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	info, err := a.client.Enqueue(task)
	if err != nil {
		logger.Error(ctx, "推送"+task.Type()+"队列失败:", err.Error())
		return nil, err
	}
	logger.Warn(ctx, "Send \"+task.Type()+\"task to queue success: ", info)
	return info, nil
}
