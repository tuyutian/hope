package providers

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hibiken/asynq"

	ossRepo "backend/internal/infras/oss"
	"backend/internal/infras/task"
)

type Option func(repos *Repositories)

func WithOssRepo(ossClient *oss.Client, bucketName string) Option {
	return func(repos *Repositories) {
		aliyunOssRepo := ossRepo.NewAliyunOSSRepository(bucketName, ossClient)
		repos.AliyunOssRepo = aliyunOssRepo
	}
}
func WithAsynqRepo(asynqClient *asynq.Client) Option {
	return func(repos *Repositories) {
		asynqRepo := task.NewAsynqRepository(asynqClient)
		repos.AsyncRepo = asynqRepo
	}
}
