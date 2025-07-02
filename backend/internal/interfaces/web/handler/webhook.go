package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/application/apps"
	"backend/internal/application/orders"
	"backend/internal/application/products"
	"backend/internal/application/users"
	orderEntity "backend/internal/domain/entity/orders"
	productEntity "backend/internal/domain/entity/products"
	"backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/response/message"
	"backend/pkg/utils"
)

type WebHookHandler struct {
	response.BaseHandler
	productService *products.ProductService
	userService    *users.UserService
	orderService   *orders.OrderService
	appService     *apps.AppService
}

func NewWebHookHandler(services *application.Services) *WebHookHandler {
	return &WebHookHandler{productService: services.ProductService, userService: services.UserService, orderService: services.OrderService, appService: services.AppService}
}

func (w *WebHookHandler) Uninstall(ctx *gin.Context) {
	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")
	appId := "insurance"
	if shopDomain != "" {
		ctxWithTrace := ctx.Request.Context()
		_ = w.userService.Uninstall(ctxWithTrace, appId, shopDomain)
	}
}

type WebhookData struct {
	ID int64 `json:"id"`
}

func (w *WebHookHandler) OrderUpdate(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()

	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookData WebhookData

	if err := json.Unmarshal(bodyBytes, &webhookData); err != nil {
		logger.Error(ctxWithTrace, "订单删除解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	_ = w.orderService.OrderSync(ctxWithTrace, orderEntity.OrderWebHookReq{
		Shop:    shopDomain,
		OrderId: webhookData.ID,
	})

	w.Success(ctx, "", nil)
}

func (w *WebHookHandler) OrderDel(ctx *gin.Context) {
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

func (w *WebHookHandler) ProductUpdate(ctx *gin.Context) {
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

func (w *WebHookHandler) ProductDel(ctx *gin.Context) {
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

type WebhookShopData struct {
	PlanDisplayName string `json:"plan_display_name"`
}

func (w *WebHookHandler) ShopUpdate(ctx *gin.Context) {
	ctxWithTrace := ctx.Request.Context()
	shopDomain := ctx.GetHeader("X-Shopify-Shop-Domain")

	bodyBytes, _ := io.ReadAll(ctx.Request.Body)

	var webhookShopData WebhookShopData

	if err := json.Unmarshal(bodyBytes, &webhookShopData); err != nil {
		logger.Error(ctxWithTrace, "ShopUpdate解析JSON失败", err, "body:", string(bodyBytes))
		utils.CallWilding(err.Error())
		return
	}

	err := w.userService.SyncShopifyUserInfo(ctxWithTrace, shopDomain, webhookShopData.PlanDisplayName)
	if err != nil {
		utils.CallWilding(err.Error())
		return
	}

	w.Success(ctx, "", nil)
}

func (w *WebHookHandler) Customers(ctx *gin.Context) {
	hmacHeader := ctx.GetHeader("X-Shopify-Hmac-Sha256")
	appId := ctx.Param("appId")
	ctxWithTrace := ctx.Request.Context()
	appConfig, err := w.appService.GetAppConfig(ctxWithTrace, appId)
	if err != nil {
		utils.CallWilding(appId + " config get with error" + err.Error())
	}
	if hmacHeader == "" {
		w.Fail(ctx, http.StatusUnauthorized, message.ErrorUnauthorized.Error(), nil)
		return
	}

	// 读取请求体数据
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		w.Fail(ctx, http.StatusBadGateway, message.ErrorBadRequest.Error(), nil)
		return
	} else {
		// 验证 HMAC
		if !utils.VerifyWebhook(data, hmacHeader, appConfig.ApiSecret) {
			w.Fail(ctx, http.StatusBadRequest, message.ErrorBadRequest.Error(), nil)
			return
		}
	}
	w.Success(ctx, "", nil)

}
