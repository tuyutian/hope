package files

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	shopifyEntity "backend/internal/domain/entity/shopifys"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/ctxkeys"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type FileService struct {
	productGraphqlRepo shopifyRepo.ProductGraphqlRepository
}

func NewFileService(repos *providers.Repositories) *FileService {
	return &FileService{
		productGraphqlRepo: repos.ProductGraphqlRepo,
	}
}

// UploadProductImageToShopify 上传文件到 Shopify
func (s *FileService) UploadProductImageToShopify(ctx context.Context, fileHeader *multipart.FileHeader, altText string) (*shopifyEntity.ImageMedia, error) {
	// 从上下文获取 Shopify GraphQL 客户端
	client := ctx.Value(ctxkeys.ShopifyGraphqlClient).(*shopify_graphql.GraphqlClient)
	// 1. 获取文件基本信息
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	contentType := fileHeader.Header.Get("Content-Type")
	fileName := utils.Uuid()
	// 2. 创建 staged upload
	stagedInput := shopifyEntity.StagedUploadInput{
		Filename:   fileName,
		MimeType:   contentType,
		Resource:   "IMAGE", // 根据实际用途选择
		HttpMethod: "POST",
	}
	fmt.Println("stagedInput:", zap.Any("stagedInput", stagedInput))
	s.productGraphqlRepo.WithClient(client)

	stagedTargets, err := s.productGraphqlRepo.StagedUploadsCreate(ctx, stagedInput)
	if err != nil {
		return nil, fmt.Errorf("创建 staged upload 失败: %w", err)
	}
	if stagedTargets == nil {
		return nil, fmt.Errorf("stagedTargets is nil")
	}
	// 3. 上传文件到预签名 URL
	stagedTarget := (*stagedTargets)[0]
	location, err := uploadToSignedURL(ctx, stagedTarget, fileContent)
	if err != nil {
		return nil, fmt.Errorf("上传到临时存储失败: %w", err)
	}
	originSource := stagedTarget.ResourceURL
	if location != "" {
		originSource = location
	}
	fileInput := shopifyEntity.FileCreateInput{
		Alt:         altText,
		ContentType: shopifyEntity.ImageType,
		// 使用返回的资源 URL
		OriginalSource: originSource,
	}

	files, err := s.productGraphqlRepo.FileCreate(ctx, fileInput)
	if err != nil {
		return nil, fmt.Errorf("创建文件记录失败: %w", err)
	}
	fmt.Println("files:", zap.Any("files", files))
	// 返回文件 ID 或预览 URL
	if len(*files) > 0 {
		uploadFile := (*files)[0]
		// 上传成功
		if uploadFile.FileStatus == "UPLOADED" {
			mediaImage, err := s.productGraphqlRepo.GetImageMedia(ctx, uploadFile.ID)
			if err != nil {
				return nil, fmt.Errorf("can't get image media: %s", err.Error())
			}
			for mediaImage.FileStatus == "PROCESSING" {
				time.Sleep(1 * time.Second)
				mediaImage, err = s.productGraphqlRepo.GetImageMedia(ctx, uploadFile.ID)
				if err != nil {
					return nil, fmt.Errorf("can't get image media: %s", err.Error())
				}
			}
			return mediaImage, nil
		}
	}
	return nil, fmt.Errorf("未获取到上传文件信息")
}

// uploadToSignedURL 上传文件到预签名 URL
func uploadToSignedURL(ctx context.Context, target shopifyEntity.StagedTarget, fileContent []byte) (string, error) {
	// 创建 multipart form 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加所有参数
	for _, param := range target.Parameters {
		err := writer.WriteField(param.Name, param.Value)
		if err != nil {
			return "", fmt.Errorf("写入参数失败: %w", err)
		}
	}

	// 添加文件内容
	part, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return "", fmt.Errorf("创建文件表单失败: %w", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return "", fmt.Errorf("写入文件内容失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭 writer 失败: %w", err)
	}

	// 发送请求
	req, err := http.NewRequest("POST", target.URL, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 如果是 XML 响应，解析响应体
	if strings.Contains(resp.Header.Get("Content-Type"), "application/xml") {
		var result struct {
			Location string `xml:"Location"`
			Bucket   string `xml:"Bucket"`
			Key      string `xml:"Key"`
		}

		if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
			logger.Error(ctx, "解析 XML 响应失败", "error", err)
			return "", err
		} else {
			logger.Info(ctx, "文件上传结果",
				"location", result.Location,
				"bucket", result.Bucket,
				"key", result.Key,
			)
			return result.Location, nil
		}
	}

	// 如果要读取响应体用于其他目的，需要重新读取
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 打印完整响应信息
	logger.Info(ctx, "上传响应",
		"statusCode", resp.StatusCode,
		"headers", resp.Header,
		"body", string(respBody),
	)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("上传失败，状态码: %d，响应内容: %s", resp.StatusCode, string(respBody))
	}

	return "", nil
}
