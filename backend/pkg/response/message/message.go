package message

import "errors"

var (
	ErrorAccountError = errors.New("账号被锁定")
	ErrorExist        = errors.New("已存在")
	ErrorUnknow       = errors.New("未知错误")
	ErrorBadRequest   = errors.New("参数错误")
	ErrorNeedLogin    = errors.New("用户未登录")
	// ErrInvalidAccount 无效账户
	ErrInvalidAccount = errors.New("invalid account")
	ErrUploadFailed   = errors.New("file upload failed")
	ErrorUnauthorized = errors.New("unauthorized")
)
