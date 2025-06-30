package shopify_graphql

type GraphqlOption func(*GraphqlClient)

func WithShop(shop string) GraphqlOption {
	return func(client *GraphqlClient) {
		client.shopName = shop
	}
}
func WithAccessToken(accessToken string) GraphqlOption {
	return func(client *GraphqlClient) {
		client.accessToken = accessToken
	}
}

func WithVersion(version string) GraphqlOption {
	return func(client *GraphqlClient) {
		client.version = version
	}
}

func WithApiPathPrefix(apiPathPrefix string) GraphqlOption {
	return func(client *GraphqlClient) {
		client.apiPathPrefix = apiPathPrefix
	}
}
