package shops

import (
	"context"
	"fmt"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
)

type ThemeGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

var _ shopifys.ThemeGraphqlRepository = (*ThemeGraphqlRepoImpl)(nil)

func NewThemeGraphqlRepository() shopifys.ThemeGraphqlRepository {
	return &ThemeGraphqlRepoImpl{}
}

func (t *ThemeGraphqlRepoImpl) GetMainThemeSettingJson(ctx context.Context) (string, error) {

	query := `query GetMainThemeSettingJson($filenames: [String!]!, $roles: [ThemeRole!]!) {
  themes(first: 1, roles: $roles) {
    nodes {
      files(filenames: $filenames) {
        nodes {
          body {
            ... on OnlineStoreThemeFileBodyBase64 {
              contentBase64
            }
            ... on OnlineStoreThemeFileBodyText {
              content
            }
            ... on OnlineStoreThemeFileBodyUrl {
              url
            }
          }
          checksumMd5
          contentType
          createdAt
          filename
          size
          updatedAt
        }
        userErrors {
          code
          filename
        }
      }
      createdAt
      id
      name
      prefix
      processing
      processingFailed
      role
      themeStoreId
      updatedAt
    }
  }
}`
	variables := map[string]interface{}{
		"filenames": []string{"config/settings_data.json"},
		"roles":     []string{"MAIN"},
	}

	var response struct {
		Themes struct {
			Nodes []shopifyEntity.OnlineStoreTheme `json:"nodes"`
		} `json:"themes"`
	}
	err := t.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return "", fmt.Errorf("查询店铺主题设置信息失败: %w", err)
	}
	if response.Themes.Nodes != nil && len(response.Themes.Nodes) > 0 && response.Themes.Nodes[0].Files.Nodes != nil && len(response.Themes.Nodes[0].Files.Nodes) > 0 {
		return response.Themes.Nodes[0].Files.Nodes[0].Body.Content, nil

	}
	return "", nil
}
