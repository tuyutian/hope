// Package products graphql 接口基础设施
package products

import (
	"context"
	"fmt"

	productEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
	"backend/pkg/utils"
)

var _ shopifys.ProductGraphqlRepository = (*productGraphqlRepoImpl)(nil)

type productGraphqlRepoImpl struct {
	shopify_graphql.Graphql
}

func NewProductGraphqlRepository() shopifys.ProductGraphqlRepository {
	return &productGraphqlRepoImpl{}
}

// CreateProduct 创建产品
func (c *productGraphqlRepoImpl) CreateProductWithMedia(ctx context.Context, productInput productEntity.ProductCreateInput, mediaInput []productEntity.CreateMediaInput) (*productEntity.ProductCreateResponse, error) {
	mutation := `
        mutation productCreate($product: ProductCreateInput!, $media: [CreateMediaInput!]) {
            productCreate(product: $product, media: $media) {
                product {
                    id
                    title
                    handle
                    status
                    createdAt
                    media(first: 10) {
						nodes {
							alt
							mediaContentType
							preview {
								status
								image {
									url
									id
								}
							}
						}
                    }
                    variants(first: 10) {
						nodes {
							id
							sku
							price
						}
                    }
                }
                userErrors {
                    field
                    message
                }
            }
        }
    `

	variables := map[string]interface{}{
		"product": productInput,
		"media":   mediaInput,
	}

	var response productEntity.ProductCreateResponse
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("创建产品失败: %w", err)
	}

	// 检查用户错误
	if len(response.ProductCreate.UserErrors) > 0 {
		return nil, fmt.Errorf("创建产品错误: %s", response.ProductCreate.UserErrors[0].Message)
	}

	return &response, nil
}

// GetProduct 查询产品
func (c *productGraphqlRepoImpl) GetProduct(ctx context.Context, productID string) (*productEntity.ProductResponse, error) {
	query := `
        query getProduct($id: ID!) {
            product(id: $id) {
                id
                title
                handle
                status
                descriptionHtml
                productType
                vendor
                tags
                createdAt
                updatedAt
                images(first: 10) {
                    edges {
                        node {
                            id
                            url
                            altText
                        }
                    }
                }
                variants(first: 50) {
                    edges {
                        node {
                            id
                            title
                            sku
                            price
                            compareAtPrice
                            inventoryQuantity
                        }
                    }
                }
            }
        }
    `

	variables := map[string]interface{}{
		"id": productID,
	}

	var response productEntity.ProductResponse
	err := c.Client.Query(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("查询产品失败: %w", err)
	}

	return &response, nil
}

// AddProductImages 为产品添加图片
func (c *productGraphqlRepoImpl) AddProductImages(ctx context.Context, productID string, images []productEntity.ProductImageInput) error {
	mutation := `
        mutation productImageUpdate($productId: ID!, $images: [ImageInput!]!) {
            productImageUpdate(productId: $productId, images: $images) {
                product {
                    id
                    images(first: 20) {
                        edges {
                            node {
                                id
                                url
                                altText
                            }
                        }
                    }
                }
                userErrors {
                    field
                    message
                }
            }
        }
    `

	variables := map[string]interface{}{
		"productId": productID,
		"images":    images,
	}

	var response struct {
		ProductImageUpdate struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"productImageUpdate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return fmt.Errorf("添加产品图片失败: %w", err)
	}

	if len(response.ProductImageUpdate.UserErrors) > 0 {
		return fmt.Errorf("添加产品图片错误: %s", response.ProductImageUpdate.UserErrors[0].Message)
	}

	return nil
}

func (c *productGraphqlRepoImpl) DeleteVariant(ctx context.Context, productID int64, variantID int64) error {
	productGid := fmt.Sprintf("gid://shopify/Product/%d", productID)
	variantGid := fmt.Sprintf("gid://shopify/ProductVariant/%d", variantID)

	mutation := `
		mutation bulkDeleteProductVariants($productId: ID!, $variantsIds: [ID!]!) {
		  productVariantsBulkDelete(productId: $productId, variantsIds: $variantsIds) {
			product {
			  id
			}
			userErrors {
			  field
			  message
			}
		  }
		}
	`

	variables := map[string]interface{}{
		"productId":   productGid,
		"variantsIds": []string{variantGid},
	}
	var response struct {
		ProductVariantsBulkDelete struct {
			Product struct {
				ID string `json:"id"`
			} `json:"product"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"productVariantsBulkDelete"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return err
	}

	if len(response.ProductVariantsBulkDelete.UserErrors) > 0 {
		return fmt.Errorf("error get product: %v", response.ProductVariantsBulkDelete.UserErrors[0].Message)
	}
	return nil
}

func (c *productGraphqlRepoImpl) CreateVariants(ctx context.Context, productID int64, input []*productEntity.VariantCreateInput) ([]map[string]interface{}, error) {
	productGid := fmt.Sprintf("gid://shopify/Product/%d", productID)
	mutation := `
		mutation ProductVariantsCreate($productId: ID!, $variants: [ProductVariantsBulkInput!]!) {
		  productVariantsBulkCreate(productId: $productId, variants: $variants) {
			productVariants {
			    id
				inventoryItem {
                    id
                    sku
	            }
			}
			userErrors {
			  field
			  message
			}
		  }
		}
	`

	variables := map[string]interface{}{
		"productId": productGid,
		"variants":  input,
	}

	var response struct {
		ProductVariantsBulkCreate struct {
			ProductVariants []struct {
				ID string `json:"id"`
				//Barcode       string `json:"barcode"`
				InventoryItem struct {
					ID  string `json:"id"`
					Sku string `json:"sku"`
				} `json:"inventoryItem"`
			} `json:"productVariants"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"productVariantsBulkCreate"`
	}

	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return nil, err
	}

	if len(response.ProductVariantsBulkCreate.UserErrors) > 0 {
		return nil, fmt.Errorf("error creating variant: %v", response.ProductVariantsBulkCreate.UserErrors[0].Message)
	}

	var variants []map[string]interface{}
	for _, edge := range response.ProductVariantsBulkCreate.ProductVariants {
		variants = append(variants, map[string]interface{}{
			"id":           edge.ID,
			"inventory_id": edge.InventoryItem.ID,
			"sku":          edge.InventoryItem.Sku,
		})
	}

	return variants, nil
}

func (c *productGraphqlRepoImpl) UpdateProduct(ctx context.Context, productID int64, product productEntity.ProductUpdateInput, media []productEntity.CreateMediaInput) error {

	product.Id = fmt.Sprintf("gid://shopify/Product/%d", productID)

	mutation := `
		mutation UpdateProductWithNewMedia($product: ProductUpdateInput!, $media: [CreateMediaInput!]) {
		  productUpdate(product: $product, media: $media) {
			product {
			  id
			}
			userErrors {
			  field
			  message
			}
		  }
		}
	`
	variables := map[string]interface{}{
		"product": product,
		//"media":   media,
	}

	if len(media) > 0 {
		variables["media"] = media
	}

	var response struct {
		ProductUpdate struct {
			Product struct {
				ID string `json:"id"`
			} `json:"product"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"productUpdate"`
	}
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return err
	}

	if len(response.ProductUpdate.UserErrors) > 0 {
		return fmt.Errorf("error creating product: %v", response.ProductUpdate.UserErrors[0].Message)
	}
	return nil
}

func (c *productGraphqlRepoImpl) PublishProduct(ctx context.Context, productID int64, publicationID int64) error {
	shopifyProductID := fmt.Sprintf("gid://shopify/Product/%d", productID)
	shopifyPublicationID := fmt.Sprintf("gid://shopify/Product/%d", publicationID)

	mutation := `mutation publishablePublish($id: ID!, $input: [PublicationInput!]!) {
		publishablePublish(id: $id, input: $input) {
			publishable {
				availablePublicationsCount {
					count
				}
			}
			userErrors {
				field
				message
			}
		}
	}`

	// 设置 GraphQL 请求变量
	variables := map[string]interface{}{
		"id":    shopifyProductID,
		"input": []map[string]interface{}{{"publicationId": shopifyPublicationID}},
	}
	var response struct {
		PublishablePublish struct {
			Publishable struct {
				AvailablePublicationsCount struct {
					Count int `json:"count"`
				} `json:"availablePublicationsCount"`
			} `json:"publishable"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"publishablePublish"`
	}
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return err
	}

	if len(response.PublishablePublish.UserErrors) > 0 {
		return fmt.Errorf("error creating product: %v", response.PublishablePublish.UserErrors[0].Message)
	}
	return nil
}

func (c *productGraphqlRepoImpl) CollectionProduct(ctx context.Context, shopifyProductID int64, collectionID int64) error {

	// 拼接 Shopify GID 格式
	collectionGID := fmt.Sprintf("gid://shopify/Collection/%d", collectionID)
	shopifyProductGID := fmt.Sprintf("gid://shopify/Product/%d", shopifyProductID)

	// 定义 GraphQL Mutation
	mutation := `
		mutation collectionAddProductsV2($id: ID!, $productIds: [ID!]!) {
			collectionAddProductsV2(id: $id, productIds: $productIds) {
				userErrors {
					field
					message
				}
			}
		}
	`

	// 设置 GraphQL 请求变量
	variables := map[string]interface{}{
		"id":         collectionGID,
		"productIds": []string{shopifyProductGID},
	}
	var response struct {
		CollectionAddProductsV2 struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"collectionAddProductsV2"`
	}
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return err
	}

	if len(response.CollectionAddProductsV2.UserErrors) > 0 {
		return fmt.Errorf("error creating product: %v", response.CollectionAddProductsV2.UserErrors[0].Message)
	}
	return nil
}

func (c *productGraphqlRepoImpl) UpdateVariants(ctx context.Context, productId int64, variants []*productEntity.VariantUpdateInput) error {

	productGid := fmt.Sprintf("gid://shopify/Product/%d", productId)

	mutation := `
		 mutation productVariantsBulkUpdate($productId: ID!, $variants: [ProductVariantsBulkInput!]!) {
            productVariantsBulkUpdate(productId: $productId, variants: $variants) {
                product {
                    id
                }
                userErrors {
                    field
                    message
                }
            }
        }
	`
	variables := map[string]interface{}{
		"productId": productGid,
		"variants":  variants,
	}

	var response struct {
		ProductVariantsBulkUpdate struct {
			Product struct {
				ID string `json:"id"`
			} `json:"product"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"productVariantsBulkUpdate"`
	}
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		return err
	}

	if len(response.ProductVariantsBulkUpdate.UserErrors) > 0 {
		return fmt.Errorf("error creating product: %v", response.ProductVariantsBulkUpdate.UserErrors[0].Message)
	}
	return nil
}

func (c *productGraphqlRepoImpl) GetCollectionList(ctx context.Context) ([]struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}, error) {
	query := `
		query collections($first: Int!, $after: String) {
			collections(first: $first, after: $after) {
				edges {
					node {
						id
						title
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	first := 100
	var after *string
	hasNextPage := true

	var collections []struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}

	// 返回查询结果
	for hasNextPage {
		variables := map[string]interface{}{
			"first": first,
		}

		if after != nil {
			variables["after"] = *after
		}

		// 初始化返回的数据结构
		var response struct {
			Collections struct {
				Edges []struct {
					Node struct {
						ID    string `json:"id"`
						Title string `json:"title"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"collections"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		err := c.Client.Query(ctx, query, variables, &response)
		if err != nil {
			return nil, err
		}

		for _, edge := range response.Collections.Edges {
			collections = append(collections, struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			}{
				ID:    utils.GetIdFromShopifyGraphqlId(edge.Node.ID),
				Title: edge.Node.Title,
			})
		}

		hasNextPage = response.Collections.PageInfo.HasNextPage
		if hasNextPage {
			after = &response.Collections.PageInfo.EndCursor
		}
	}

	// 返回查询结果
	return collections, nil
}
