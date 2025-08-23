package shopify_graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/machinebox/graphql"

	"backend/pkg/utils"
)

var (
	defaultVersion       = "2025-07"
	defaultApiPathPrefix = "admin/api"
)

type GraphqlClient struct {
	client        *graphql.Client
	shopName      string
	accessToken   string
	version       string
	apiPathPrefix string
}

func NewGraphqlClient(shopName, accessToken string, opts ...GraphqlOption) *GraphqlClient {
	endpoint := fmt.Sprintf("https://%s.myshopify.com/%s/%s/graphql.json", shopName, defaultApiPathPrefix, defaultVersion)

	client := graphql.NewClient(endpoint)
	client.Log = func(s string) {
		if strings.Contains(s, "errors") {
			utils.CallWilding(s)
		}
		fmt.Println(s)
	}
	graphqlClient := &GraphqlClient{
		client:        client,
		shopName:      shopName,
		accessToken:   accessToken,
		version:       defaultVersion,
		apiPathPrefix: defaultApiPathPrefix,
	}
	// apply any options
	for _, opt := range opts {
		opt(graphqlClient)
	}
	return graphqlClient
}

// 设置请求头
func (c *GraphqlClient) setHeaders(req *graphql.Request) {
	// 安全检查
	if req == nil {
		fmt.Println("Warning: graphql.Request is nil")
		return
	}

	// 确保 Header 已初始化
	if req.Header == nil {
		req.Header = make(map[string][]string)
	}

	// 设置必要的请求头
	req.Header.Set("X-Shopify-Access-Token", c.accessToken)
	req.Header.Set("Content-Type", "application/json")

}

// Query 执行 GraphQL 查询
func (c *GraphqlClient) Query(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(query)

	// 添加变量
	for key, value := range variables {
		req.Var(key, value)
	}
	// 设置请求头
	c.setHeaders(req)

	return c.client.Run(ctx, req, response)
}

// Mutate 执行 GraphQL 变更
func (c *GraphqlClient) Mutate(ctx context.Context, mutation string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(mutation)

	// 添加变量
	for key, value := range variables {
		req.Var(key, value)
	}

	// 设置请求头
	c.setHeaders(req)

	return c.client.Run(ctx, req, response)
}
