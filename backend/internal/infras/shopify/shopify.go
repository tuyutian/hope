package shopify

import (
	"context"
	"net/url"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/pkg/ctxkeys"
	"backend/pkg/response/message"
)

var _ shopifyRepo.ShopifyRepository = (*shopifyRepoImpl)(nil)

type shopifyRepoImpl struct {
}

func NewShopifyRepository() shopifyRepo.ShopifyRepository {
	return &shopifyRepoImpl{}
}

var accessTokenRelPath = "admin/oauth/access_token"

func (s *shopifyRepoImpl) RequestOfflineSessionToken(ctx context.Context, token string) (*shopifyEntity.Token, error) {

	app := ctx.Value(ctxkeys.ShopifyApp).(*goshopify.App)
	client := ctx.Value(ctxkeys.ShopifyClient).(*goshopify.Client)
	data := struct {
		ClientId           string `json:"client_id"`
		ClientSecret       string `json:"client_secret"`
		GrantType          string `json:"grant_type"`
		SubjectToken       string `json:"subject_token"`
		SubjectTokenType   string `json:"subject_token_type"`
		RequestedTokenType string `json:"requested_token_type"`
	}{
		ClientId:           app.ApiKey,
		ClientSecret:       app.ApiSecret,
		GrantType:          "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectToken:       token,
		SubjectTokenType:   "urn:ietf:params:oauth:token-type:id_token",
		RequestedTokenType: "urn:shopify:params:oauth:token-type:offline-access-token",
	}
	req, err := client.NewRequest(ctx, "POST", accessTokenRelPath, data, nil)
	if err != nil {
		return nil, err
	}

	sessionToken := new(shopifyEntity.Token)
	err = client.Do(req, sessionToken)
	return sessionToken, err
}

func (s *shopifyRepoImpl) GetShopName(ctx context.Context, shopUrl string) (string, error) {
	if !strings.HasPrefix(shopUrl, "https://") {
		shopUrl = "https://" + shopUrl
	}
	parsedURL, err := url.Parse(shopUrl)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()

	if strings.HasSuffix(host, ".myshopify.com") {
		shopName := strings.TrimSuffix(host, ".myshopify.com")
		return shopName, nil
	}
	return "", message.ErrInvalidAccount
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
