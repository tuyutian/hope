package middleware

import (
	"fmt"
	"strings"

	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CspMiddleware CSP 中间件
type CspMiddleware struct {
	isEmbeddedApp bool
}

// NewCspMiddleware 创建 CSP 中间件
func NewCspMiddleware(isEmbeddedApp bool) *CspMiddleware {
	return &CspMiddleware{
		isEmbeddedApp: isEmbeddedApp,
	}
}

// Csp 设置 Content Security Policy 头
func (m *CspMiddleware) Csp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 获取 shop 域名
		shop := c.Query("shop")
		if shop == "" {
			// 尝试从用户信息中获取
			shop = c.GetHeader("X-Shopify-Shop-Domain")
		}

		// 清理 shop 域名
		shop, _ = utils.GetShopName(shop)

		// 确定允许的域名
		var allowedDomains string
		if m.isEmbeddedApp {
			if shop != "" {
				allowedDomains = fmt.Sprintf("https://%s https://admin.shopify.com", shop)
			} else {
				allowedDomains = "*.myshopify.com https://admin.shopify.com"
			}
		} else {
			allowedDomains = "'none'"
		}

		// 设置 CSP 头
		cspHeader := m.buildCspHeader(allowedDomains)
		c.Header("Content-Security-Policy", cspHeader)

		logger.Debug(ctx, "CSP 头已设置",
			"shop", shop,
			"is_embedded", m.isEmbeddedApp,
			"allowed_domains", allowedDomains)

		c.Next()
	}
}

// buildCspHeader 构建 CSP 头
func (m *CspMiddleware) buildCspHeader(frameAncestors string) string {
	// 基础 CSP 指令
	styleSrc := "style-src 'unsafe-inline' 'self' https://hope-cu0.pages.dev https://*.protectifyapp.com https://*.cloudflare.com https://*.gstatic.com https://*.googleapis.com "
	connectSrc := "connect-src 'self' https://*.protectifyapp.com https://cloudflareinsights.com/ https://*.google-analytics.com https://*.googleapis.com "
	fontSrc := "font-src 'self' https://*.gstatic.com https://*.protectifyapp.com https://*.cloudflare.com data:  "
	scriptSrc := "script-src 'self' https://hope-cu0.pages.dev https://*.protectifyapp.com https://*.cloudflare.com https://static.cloudflareinsights.com  https://*.google-analytics.com https://*.google.com https://*.googleapis.com https://www.googletagmanager.com 'unsafe-inline' 'unsafe-eval' blob:"
	imgSrc := "img-src * 'self' https://www.google-analytics.com www.googletagmanager.com data: https:"
	defaultSrc := "default-src 'self' blob:"

	// 开发环境添加本地域名
	if gin.Mode() == gin.DebugMode {
		devCsp := " ws://*:9527 http://*:9527 http://*:8080"
		styleSrc += devCsp
		connectSrc += devCsp
		fontSrc += devCsp
		scriptSrc += devCsp
		imgSrc += devCsp
	}

	// 组合所有指令
	cspParts := []string{
		styleSrc,
		connectSrc,
		fontSrc,
		scriptSrc,
		imgSrc,
		defaultSrc,
		fmt.Sprintf("frame-ancestors %s", frameAncestors),
	}

	return strings.Join(cspParts, "; ")
}
