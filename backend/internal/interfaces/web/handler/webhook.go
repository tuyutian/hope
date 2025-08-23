package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/apps"
	"backend/internal/application/orders"
	"backend/internal/application/products"
	"backend/internal/application/users"
	orderEntity "backend/internal/domain/entity/orders"
	productEntity "backend/internal/domain/entity/products"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/response/code"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

type WebHookHandler struct {
	response.BaseHandler
	productService      *products.ProductService
	userService         *users.UserService
	orderService        *orders.OrderService
	appService          *apps.AppService
	billingService      *users.BillingService
	subscriptionService *users.SubscriptionService
}

func NewWebHookHandler(services *application.Services) *WebHookHandler {
	return &WebHookHandler{
		productService: services.ProductService, userService: services.UserService, orderService: services.OrderService, appService: services.AppService,
		billingService:      services.BillingService,
		subscriptionService: services.SubscriptionService,
	}
}

type WebhookData struct {
	ID int64 `json:"id"`
}

type WebhookShopData struct {
	PlanDisplayName string `json:"plan_display_name"`
}

func (w *WebHookHandler) Shopify(ctx *gin.Context) {
	// 获取已注册的 webhook topics
	registerTopics := shopifyRepo.ShopifyWebhookTopics
	complianceTopics := shopifyRepo.ShopifyComplianceTopics
	ctxs := ctx.Request.Context()
	// 获取 Shopify 签名
	signature := ctx.GetHeader("X-Shopify-Hmac-Sha256")
	if signature == "" {
		fmt.Println("警告: 缺少 Shopify 签名头")
		w.Fail(ctx, http.StatusUnauthorized, message.ErrorUnauthorized.Error(), nil)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Printf("读取请求体失败: %v\n", err)
		w.Fail(ctx, http.StatusBadGateway, message.ErrorBadRequest.Error(), nil)
		return
	}

	// 重新设置请求体供后续使用
	ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))
	// 验证 webhook 签名
	if !w.appService.VerifyWebhook(ctxs, signature, body) {
		w.Fail(ctx, http.StatusUnauthorized, message.ErrorUnauthorized.Error(), nil)
		return
	}

	// 获取 webhook topic
	topic := ctx.GetHeader("X-Shopify-Topic")
	if topic == "" {
		w.Error(ctx, code.BadRequest, "缺少 X-Shopify-Topic 头", "")
		return
	}
	logger.Warn(ctxs, "webhook topic"+topic)
	// 验证是否为已注册的 topic
	if !w.isRegisteredTopic(topic, registerTopics) && !w.isRegisteredTopic(topic, complianceTopics) {
		w.Error(ctx, code.BadRequest, "未注册的 webhook topic", topic)
		return
	}
	appID := w.appService.GetAppID(ctx.Request.Context())
	// 根据已注册的 topic 处理不同类型的回调
	switch topic {
	case "orders/updated":
		w.handleOrderUpdated(ctx, appID, body)
	case "orders/delete":
		w.handleOrderDeleted(ctx, appID, body)
	case "app_subscriptions/update":
		w.handleAppSubscriptionUpdate(ctx, appID, body)
	case "app_subscriptions/approaching_capped_amount":
		w.handleAppSubscriptionApproachingCappedAmount(ctx, appID, body)
	case "app/uninstalled":
		w.handleAppUninstalled(ctx, appID, body)
	case "products/update":
		w.handleProductUpdated(ctx, appID, body)
	case "products/delete":
		w.handleProductDeleted(ctx, appID, body)
	case "shop/update":
		w.handleShopUpdated(ctx, appID, body)
	case "shop/redact":
		w.handleShopUpdated(ctx, appID, body)
	case "customers/redact":
		w.Customers(ctx, appID, body)
	case "customers/data_request":
		w.Customers(ctx, appID, body)
	default:
		// 记录未处理的 topic（虽然通过了注册验证，但可能是新增的）
		fmt.Printf("收到已注册但未处理的 webhook topic: %s\n", topic)
		w.Success(ctx, "webhook received but not processed", nil)
	}
}

// isRegisteredTopic 检查 topic 是否为已注册的
func (w *WebHookHandler) isRegisteredTopic(topic string, registeredTopics []string) bool {
	for _, registeredTopic := range registeredTopics {
		if topic == registeredTopic {
			return true
		}
	}
	return false
}

// handleOrderUpdated 处理订单更新
func (w *WebHookHandler) handleOrderUpdated(ctx *gin.Context, appID string, body []byte) {
	ctxWithTrace := ctx.Request.Context()

	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookData WebhookData

	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		logger.Error(ctxWithTrace, "订单删除解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	_ = w.orderService.OrderSync(ctxWithTrace, appID, orderEntity.OrderWebHookReq{
		Shop:    shopDomain,
		OrderId: webhookData.ID,
	})

	w.Success(ctx, "", nil)
}

// handleOrderDeleted 处理订单删除
func (w *WebHookHandler) handleOrderDeleted(ctx *gin.Context, appID string, body []byte) {
	ctxWithTrace := ctx.Request.Context()

	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")
	// 2. 读取 Body
	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookData WebhookData

	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		logger.Error(ctxWithTrace, "订单删除解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	_ = w.orderService.OrderDel(ctxWithTrace, orderEntity.OrderWebHookReq{
		Shop:    shopDomain,
		OrderId: webhookData.ID,
	})

	w.Success(ctx, "", nil)
}

// handleAppSubscriptionUpdate 处理应用订阅更新
func (w *WebHookHandler) handleAppSubscriptionUpdate(ctx *gin.Context, appID string, body []byte) {
	fmt.Printf("处理应用订阅更新 webhook: %s\n", string(body))

	var subscription shopifyEntity.SubscriptionWebhookPayload
	if err := json.Unmarshal(body, &subscription); err != nil {
		w.Error(ctx, code.BadRequest, "解析订阅数据失败", err.Error())
		return
	}
	user, _ := w.userService.GetUserFromShopID(ctx, utils.GetIdFromShopifyGraphqlId(subscription.AdminGraphqlApiShopId))
	// 更新订阅状态
	err := w.subscriptionService.UpdateSubscriptionStatus(
		ctx.Request.Context(),
		user,
		utils.GetIdFromShopifyGraphqlId(subscription.AdminGraphqlApiId),
		subscription.Status,
	)

	if err != nil {
		w.Error(ctx, code.ServerOperationFailed, "更新订阅状态失败", err.Error())
		return
	}

	fmt.Printf("订阅更新成功: ID=%s, Status=%s\n", subscription.AdminGraphqlApiId, subscription.Status)
	ctx.JSON(http.StatusOK, gin.H{"message": "订阅更新处理成功"})
}

// handleAppSubscriptionApproachingCappedAmount 处理订阅接近上限金额
func (w *WebHookHandler) handleAppSubscriptionApproachingCappedAmount(ctx *gin.Context, appID string, body []byte) {
	fmt.Printf("处理订阅接近上限金额 webhook: %s\n", string(body))

	var subscriptionApproaching shopifyEntity.SubscriptionApproachingPayload
	if err := json.Unmarshal(body, &subscriptionApproaching); err != nil {
		w.Error(ctx, code.BadRequest, "解析订阅数据失败", err.Error())
		return
	}

	// 发送通知或处理接近上限的逻辑
	err := w.subscriptionService.HandleApproachingCappedAmount(
		ctx.Request.Context(),
		utils.GetIdFromShopifyGraphqlId(subscriptionApproaching.AdminGraphqlApiId),
	)

	if err != nil {
		w.Error(ctx, code.ServerOperationFailed, "处理接近上限金额失败", err.Error())
		return
	}

	fmt.Printf("订阅接近上限处理成功: ID=%s\n", subscriptionApproaching.AdminGraphqlApiId)
	ctx.JSON(http.StatusOK, gin.H{"message": "订阅接近上限处理成功"})
}

// handleAppUninstalled 处理应用卸载
func (w *WebHookHandler) handleAppUninstalled(ctx *gin.Context, appID string, body []byte) {
	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")
	ctxWithTrace := ctx.Request.Context()
	appId := w.appService.GetAppID(ctxWithTrace)
	if shopDomain != "" {
		_ = w.userService.Uninstall(ctxWithTrace, appId, shopDomain)
	}
}

// handleProductUpdated 处理产品更新
func (w *WebHookHandler) handleProductUpdated(ctx *gin.Context, appID string, body []byte) {
	ctxWithTrace := ctx.Request.Context()

	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookData WebhookData

	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		logger.Error(ctxWithTrace, "产品修改解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	_ = w.productService.ProductUpdate(ctxWithTrace, productEntity.ProductWebHookReq{
		Shop:      shopDomain,
		ProductId: webhookData.ID,
	})

	w.Success(ctx, "", nil)
}

// handleProductDeleted 处理产品删除
func (w *WebHookHandler) handleProductDeleted(ctx *gin.Context, appID string, body []byte) {
	ctxWithTrace := ctx.Request.Context()
	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)
	var webhookData WebhookData

	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		logger.Error(ctxWithTrace, "产品删除解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	_ = w.productService.ProductDel(ctxWithTrace, productEntity.ProductWebHookReq{
		Shop:      shopDomain,
		ProductId: webhookData.ID,
	})

	w.Success(ctx, "", nil)
}

// handleShopUpdated 处理店铺更新
func (w *WebHookHandler) handleShopUpdated(ctx *gin.Context, appID string, body []byte) {
	ctxWithTrace := ctx.Request.Context()
	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookShopData WebhookShopData

	if err := json.Unmarshal(bodyBytes, &webhookShopData); err != nil {
		logger.Error(ctxWithTrace, "ShopUpdate解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	err := w.userService.SyncShopifyUserInfo(ctxWithTrace, appID, shopDomain, webhookShopData.PlanDisplayName)
	if err != nil {
		utils.CallWilding(err.Error())
		return
	}

	w.Success(ctx, "", nil)
}

// ChargeCallback 处理 Shopify 订阅回调
func (w *WebHookHandler) ChargeCallback(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()

	appID := w.appService.GetAppID(ctxWithTrace)

	chargeIDStr := ctx.Query("charge_id")
	userIDStr := ctx.Param("userID")
	// 记录原始回调数据用于调试
	fmt.Printf("收到订阅回调,charge id is: %s, user id is: %s\n", chargeIDStr, userIDStr)
	chargeID, err := strconv.ParseInt(chargeIDStr, 10, 64)
	// 记录原始回调数据用于调试
	userID, err1 := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || err1 != nil {
		w.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	// 解析通用订阅数据
	if chargeID <= 0 || userID <= 0 {
		w.Error(ctx, code.BadRequest, message.ErrorBadRequest.Error(), "")
		return
	}
	user, err := w.userService.GetLoginUserFromID(ctxWithTrace, userID)

	if user == nil {
		w.Error(ctx, code.Unauthorized, message.ErrInvalidAccount.Error(), "")
		return
	}
	_, err = w.subscriptionService.VerifyPayment(ctxWithTrace, user, chargeID)
	if err != nil {
		logger.Error(ctxWithTrace, "charge callback verify payment error", err.Error(), "charge id is: ", chargeID, "user id is: ", userID)
	}
	appConf, err := w.appService.GetAppConfig(ctxWithTrace, appID)
	if err != nil {
		w.Error(ctx, code.NotFound, message.ErrInvalidAccount.Error(), appID)
		return
	}
	appLink := appConf.AppLink
	shopName, _ := utils.GetShopName(user.Shop)
	redirectUrl := fmt.Sprintf("https://admin.shopify.com/store/%s/apps/%s/cart", shopName, appLink)
	ctx.Redirect(http.StatusFound, redirectUrl)
}

func (w *WebHookHandler) Customers(ctx *gin.Context, appID string, body []byte) {
	w.Success(ctx, "", nil)

}
