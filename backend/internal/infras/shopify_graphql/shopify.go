package shopify_graphql

import (
	"context"

	"go.uber.org/zap"

	"backend/pkg/logger"
)

type Graphql struct {
	Client *GraphqlClient
}

func (b *Graphql) WithClient(client *GraphqlClient) {
	b.Client = client
}

func (b *Graphql) GetByID(ctx context.Context, id string, query string, response interface{}) error {

	variables := map[string]interface{}{
		"id": id,
	}
	err := b.Client.Query(ctx, query, variables, response)
	if err != nil {
		logger.Error(ctx, "Get data by ID error: "+err.Error(), zap.String("id", id), zap.Any("response", response))
		return err
	}
	return nil
}
