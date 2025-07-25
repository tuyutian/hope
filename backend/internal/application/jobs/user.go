package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"backend/internal/domain/entity/jobs"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type UserService struct {
	userRepo        users.UserRepository
	shopifyRepo     shopifyRepo.ShopifyRepository
	shopGraphqlRepo shopifyRepo.ShopGraphqlRepository
}

func NewUserService(repos *providers.Repositories) *UserService {
	return &UserService{
		userRepo:        repos.UserRepo,
		shopifyRepo:     repos.ShopifyRepo,
		shopGraphqlRepo: repos.ShopGraphqlRepo,
	}
}

func (u *UserService) HandleInitUser(ctx context.Context, t *asynq.Task) error {
	var payload jobs.InitUserPayload

	logger.Info(ctx, "init_user_queue:", "我正在消费初始化用户")
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return u.fail(ctx, 0, "payload 反序列化失败", err)
	}

	uid := payload.UserID
	user, err := u.userRepo.Get(ctx, uid)
	if err != nil || user == nil {
		return u.fail(ctx, 0, "查询用户信息错误", err)
	}

	// 初始化 Shopify client
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	u.shopGraphqlRepo.WithClient(client)
	webhooks, _ := u.shopGraphqlRepo.QueryWebhookSubscriptions(ctx, "")

	topics := shopifyRepo.ShopifyWebhookTopics
	webhookUrl := u.shopifyRepo.GetWebhookUrl(user.AppId)
	for _, topic := range topics {
		hasWebhook := false
		if webhooks != nil {
			for _, webhook := range webhooks {
				if webhook.Topic == topic {
					logger.Info(ctx, fmt.Sprintf("Update %s webhook for shop %s", topic, user.Shop))
					_ = u.shopGraphqlRepo.UpdateWebhookSubscription(ctx, webhook.Id, webhookUrl)
					hasWebhook = true
					logger.Warn(ctx, fmt.Sprintf("Update %s webhook for shop %s", topic, user.Shop),
						map[string]interface{}{
							"old": webhook.Endpoint,
							"new": webhookUrl,
						})
					break
				}
			}
		}
		if !hasWebhook {
			logger.Info(ctx, fmt.Sprintf("Register %s webhook for shop %s", topic, user.Shop))
			_ = u.shopGraphqlRepo.CreateWebhookSubscription(ctx, topic, webhookUrl)
			logger.Warn(ctx, fmt.Sprintf("Register %s webhook for shop %s", topic, user.Shop),
				map[string]interface{}{"webhook": webhookUrl})
		}
	}
	publishId, err := u.shopGraphqlRepo.GetPublicationID(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("get publication id error:%s", err.Error()))
	}
	user.PublishId = utils.GetIdFromShopifyGraphqlId(publishId)
	if err := u.userRepo.Update(ctx, user); err != nil {
		logger.Error(ctx, fmt.Sprintf("update user error:%s", err.Error()))
	}
	return nil
}

func (u *UserService) fail(ctx context.Context, uid int64, msg string, err error) error {
	logger.Error(ctx, fmt.Sprintf("init_user_queue:%d %s: %v", uid, msg, err))
	return nil
}
