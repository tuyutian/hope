package shopify

import (
	"strings"
)

// ExtractCurrencySymbol 从 moneyFormat 提取货币符号
func ExtractCurrencySymbol(moneyFormat string) string {
	// 以 "{{" 分割，获取前面的符号部分
	parts := strings.Split(moneyFormat, "{{")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// GetShopifyGraphqlId /**
func GetShopifyGraphqlId(id string) string {
	// 找到最后一个 "/" 的位置
	index := strings.LastIndex(id, "/")
	if index != -1 {
		// 从 "/" 后面开始截取字符串
		return id[index+1:]
	}
	return ""
}
