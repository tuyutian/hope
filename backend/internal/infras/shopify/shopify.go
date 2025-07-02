package shopify

import (
	"fmt"
	"strings"

	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/infras/config"
)

var _ shopifyRepo.ShopifyRepository = (*shopifyRepoImpl)(nil)

type shopifyRepoImpl struct {
	webhookHost string
}

var WebhookShopifyEndpoint = "api/v1/shopify/webhook"

func NewShopifyRepository(shopifyConf *config.Shopify) shopifyRepo.ShopifyRepository {
	return &shopifyRepoImpl{
		webhookHost: shopifyConf.WebhookHost,
	}
}
func (s *shopifyRepoImpl) GetWebhookUrl() string {

	return fmt.Sprintf("https://%s/%s", s.webhookHost, WebhookShopifyEndpoint)
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
