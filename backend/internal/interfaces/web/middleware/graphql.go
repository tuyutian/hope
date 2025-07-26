package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"backend/internal/application/users"
	"backend/internal/domain/repo"
	appRepo "backend/internal/domain/repo/apps"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/utils"
)

type ShopifyGraphqlWare struct {
	appID       string
	userService *users.UserService
	cacheRepo   repo.CacheRepository
	appRepo     appRepo.AppRepository
	shopifyRepo shopifyRepo.ShopifyRepository
}

func NewShopifyGraphqlWare(repos *providers.Repositories, userService *users.UserService) *ShopifyGraphqlWare {
	return &ShopifyGraphqlWare{
		cacheRepo:   repos.CacheRepo,
		appRepo:     repos.AppRepo,
		userService: userService,
	}
}

func (w *ShopifyGraphqlWare) ShopifyGraphqlClient() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims := w.userService.GetClaims(ctx)
		user, _ := w.userService.GetLoginUserFromID(ctx, claims.UserID)
		accessToken := "fooShop"
		if user != nil {
			accessToken = user.AccessToken
		}
		shopName, _ := utils.GetShopName(claims.Dest)
		client := shopify_graphql.NewGraphqlClient(shopName, accessToken)
		ctx = context.WithValue(ctx, ctxkeys.ShopifyGraphqlClient, client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
