package shopify

import (
	"fmt"
	"strconv"
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

// GetIdFromShopifyGraphqlId /**
func (s *shopifyRepoImpl) GetIdFromShopifyGraphqlId(gid string) int64 {
	if gid == "" {
		return 0
	}

	var idStr string
	if strings.HasPrefix(gid, "gid://shopify/") {
		parts := strings.Split(gid, "/")
		if len(parts) > 0 {
			idStr = parts[len(parts)-1]
		}
	} else {
		idStr = gid
	}

	// 将字符串转换为 int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0 // 转换失败时返回 0
	}

	return id
}
