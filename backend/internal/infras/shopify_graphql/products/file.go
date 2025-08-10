package products

import (
	"context"
	"errors"

	"go.uber.org/zap"

	productEntity "backend/internal/domain/entity/shopifys"
	"backend/pkg/logger"
)

func (c *productGraphqlRepoImpl) FileCreate(ctx context.Context, input productEntity.FileCreateInput) (*[]productEntity.FileCreated, error) {
	// 执行 GraphQL 请求
	mutation := `
mutation fileCreate($files: [FileCreateInput!]!) {
	fileCreate(files: $files) {
		files {
			id
			fileStatus
			alt
			createdAt
			fileErrors{
				code
				message
				details
			}
		}
		userErrors {
			field
			message
		}
	}
}`
	variables := map[string]interface{}{
		"files": []productEntity.FileCreateInput{
			input,
		},
	}
	var response productEntity.FileCreateResponse
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		logger.Error(ctx, "fileCreate error: "+err.Error(), zap.Any("response", response))
		return nil, err
	}
	if len(response.FileCreate.UserErrors) > 0 {
		logger.Error(ctx, "fileCreate error: "+response.FileCreate.UserErrors[0].Message, zap.Any("response", response))
		return nil, errors.New(response.FileCreate.UserErrors[0].Message)
	}
	logger.Info(ctx, "fileCreate success", zap.Any("response", response))

	return &response.FileCreate.Files, nil
}

func (c *productGraphqlRepoImpl) StagedUploadsCreate(ctx context.Context, input productEntity.StagedUploadInput) (*[]productEntity.StagedTarget, error) {
	mutation := `
mutation stagedUploadsCreate($input: [StagedUploadInput!]!) {
  stagedUploadsCreate(input: $input) {
    stagedTargets {
      url
      resourceUrl
      parameters {
        name
        value
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
		"input": []productEntity.StagedUploadInput{
			input,
		},
	}

	var response productEntity.StagedUploadsCreateResponse
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		logger.Error(ctx, "fileCreate error: "+err.Error(), zap.Any("response", response))
		return nil, err
	}
	if len(response.StagedUploadsCreate.UserErrors) > 0 {
		logger.Error(ctx, "fileCreate error: "+response.StagedUploadsCreate.UserErrors[0].Message, zap.Any("response", response))
		return nil, errors.New(response.StagedUploadsCreate.UserErrors[0].Message)
	}

	return &response.StagedUploadsCreate.StagedTargets, nil
}

func (c *productGraphqlRepoImpl) GetImageMedia(ctx context.Context, id string) (*productEntity.ImageMedia, error) {
	query := `
query GetMediaByID($id:ID!) {
  node(id: $id) {
    id
    ... on MediaImage {
      fileErrors{
        code
        details
        message
      }
      alt
      originalSource {
        url
        fileSize
      }
      fileStatus
      image {
        url
      }
    }
    ... on File {
      fileErrors{
        code
        details
        message
      }
      alt
    	preview{
        image{
          url
        }
      }
      fileStatus
    }
    
  }
}`
	var response struct {
		Node *productEntity.ImageMedia
	}
	err := c.Graphql.GetByID(ctx, id, query, &response)
	if err != nil {
		return nil, err
	}
	logger.Warn(ctx, "GetImageMedia", zap.Any("response", response.Node))
	if response.Node != nil {
		return response.Node, nil
	}
	return nil, errors.New("file not found")
}

func (c *productGraphqlRepoImpl) FileUpdate(ctx context.Context, input productEntity.FileUpdateInput) (*[]productEntity.FileUpdated, error) {
	// 执行 GraphQL 请求
	mutation := `
mutation fileUpdate($files: [FileUpdateInput!]!) {
  fileUpdate(files: $files) {
    files {
		id
		fileStatus
		alt
		createdAt
		preview {
			image {
				id
				url
				altText
				}
			status
		}
		fileErrors{
			code
			message
			details
		}
	}
    userErrors {
      field
      message
      code
    }
  }
}`
	variables := map[string]interface{}{
		"files": []productEntity.FileUpdateInput{
			input,
		},
	}
	var response productEntity.FileUpdateResponse
	err := c.Client.Mutate(ctx, mutation, variables, &response)
	if err != nil {
		logger.Error(ctx, "fileUpdate error: "+err.Error(), zap.Any("response", response))
		return nil, err
	}
	if len(response.FileUpdate.UserErrors) > 0 {
		logger.Error(ctx, "fileUpdated error: "+response.FileUpdate.UserErrors[0].Message, zap.Any("response", response))
		return nil, errors.New(response.FileUpdate.UserErrors[0].Message)
	}
	logger.Info(ctx, "fileUpdated success", zap.Any("response", response))

	return &response.FileUpdate.Files, nil
}
