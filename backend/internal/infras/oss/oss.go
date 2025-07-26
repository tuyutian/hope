package oss

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"backend/internal/domain/repo"
	"backend/pkg/logger"
)

var _ repo.AliyunOSSRepository = (*aliYunOssRepoImpl)(nil)

type aliYunOssRepoImpl struct {
	bucketName string
	client     *oss.Client
}

func NewAliyunOSSRepository(bucketName string, client *oss.Client) repo.AliyunOSSRepository {
	return &aliYunOssRepoImpl{
		bucketName: bucketName,
		client:     client,
	}
}

func (a *aliYunOssRepoImpl) UploadFile(ctx context.Context, fileName string, file *multipart.FileHeader) (string, error) {

	bucket, err := a.client.Bucket(a.bucketName)
	if err != nil {
		return "", err
	}

	fileData, err := file.Open()
	defer fileData.Close()

	err = bucket.PutObject(fileName, fileData)
	if err != nil {
		return "", err
	}
	imagePath := "https://" + a.bucketName + "." + a.client.Config.Endpoint + "/" + fileName
	fmt.Println("文件上传到：", imagePath)
	logger.Warn(ctx, "文件上传到："+imagePath)
	return imagePath, nil
}
