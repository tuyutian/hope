package repo

import (
	"context"
	"mime/multipart"
)

type AliyunOSSRepository interface {
	UploadFile(ctx context.Context, fileName string, file *multipart.FileHeader) (string, error)
}
