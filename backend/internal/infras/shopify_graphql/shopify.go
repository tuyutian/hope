package shopify_graphql

type Graphql struct {
	Client *GraphqlClient
}

func (b *Graphql) WithClient(client *GraphqlClient) {
	b.Client = client
}
