package shopify_graphql

type GraphqlOption func(g *GraphqlClient)

func WithShop(shop string) GraphqlOption {
	return func(g *GraphqlClient) {
		g.shopName = shop
	}
}
func WithAccessToken(accessToken string) GraphqlOption {
	return func(g *GraphqlClient) {
		g.accessToken = accessToken
	}
}

func WithVersion(version string) GraphqlOption {
	return func(g *GraphqlClient) {
		g.version = version
	}
}

func WithApiPathPrefix(apiPathPrefix string) GraphqlOption {
	return func(g *GraphqlClient) {
		g.apiPathPrefix = apiPathPrefix
	}
}
