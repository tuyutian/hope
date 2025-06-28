package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	goshopify "github.com/bold-commerce/go-shopify/v4"

	"backend/internal/application/users"
	"backend/internal/domain/repo"
	appRepo "backend/internal/domain/repo/apps"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/response/code"
)

type GoShopifyWare struct {
	appID       string
	userService users.UserService
	cacheRepo   repo.CacheRepository
	appRepo     appRepo.AppRepository
	shopifyRepo shopifyRepo.ShopifyRepository
}

func (w *GoShopifyWare) NewGoShopifyWare(appId string, repos *providers.Repositories, userService users.UserService) *GoShopifyWare {
	return &GoShopifyWare{
		cacheRepo:   repos.CacheRepo,
		appRepo:     repos.AppRepo,
		userService: userService,
	}
}

func (w *GoShopifyWare) GoShopify(appID string) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		shopify, err := w.appRepo.GetGoShopifyByAppID(ctx, appID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    code.AppNotfound,
				"message": http.StatusText(http.StatusNotFound),
			})
			return
		}

		claims := w.userService.GetClaims(ctx)
		user, _ := w.userService.GetLoginUserFromID(ctx, claims.UserID)
		accessToken := "fooShop"
		if user != nil {
			accessToken = user.AccessToken
		}
		shopName, _ := w.shopifyRepo.GetShopName(ctx, claims.Dest)
		client, _ := shopify.NewClient(shopName, accessToken, goshopify.WithVersion("2025-07"), goshopify.WithRetry(3))
		ctx = context.WithValue(ctx, ctxkeys.ShopifyClient, client)
		ctx = context.WithValue(ctx, ctxkeys.ShopifyApp, shopify)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
