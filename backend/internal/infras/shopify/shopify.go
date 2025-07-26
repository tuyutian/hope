package shopify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/infras/config"
)

var _ shopifyRepo.ShopifyRepository = (*shopifyRepoImpl)(nil)

type shopifyRepoImpl struct {
	webhookHost string
}

var WebhookShopifyEndpoint = "api/v1/webhook/shopify"
var PaymentCallback = "api/v1/webhook/charge_callback"

func NewShopifyRepository(shopifyConf *config.Shopify) shopifyRepo.ShopifyRepository {
	return &shopifyRepoImpl{
		webhookHost: shopifyConf.WebhookHost,
	}
}
func (s *shopifyRepoImpl) GetWebhookUrl(appID string) string {

	return fmt.Sprintf("https://%s/%s/%s", s.webhookHost, appID, WebhookShopifyEndpoint)
}
func (s *shopifyRepoImpl) GetReturnUrl(appID string, userID int64) string {
	return fmt.Sprintf("https://%s/%s/%s/%d", s.webhookHost, appID, PaymentCallback, userID)
}

// ExtractCurrencySymbol 从 moneyFormat 提取货币符号
func (s *shopifyRepoImpl) ExtractCurrencySymbol(moneyFormat string) string {
	// 以 "{{" 分割，获取前面的符号部分
	parts := strings.Split(moneyFormat, "{{")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// VerifyWebhook 验证 webhook 签名
func (s *shopifyRepoImpl) VerifyWebhook(ctx context.Context, appSecret string, signature string, body []byte) bool {

	// 计算期望的签名
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 比较签名
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
