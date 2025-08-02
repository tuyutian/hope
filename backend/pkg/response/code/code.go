package code

// HTTP状态码
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// 业务错误码（用于响应体中的code字段）
// 1级公式错误码 1～9999
// 业务模块错误码 100xxxx
// 其他模块 101xxxx 102xxxx 以此类推
var (
	Success               = 0
	BadRequest            = 1010001
	NotFound              = 1010002
	Unauthorized          = 1010003
	AppNotfound           = 1010004
	UserNotfound          = 1010005
	ThirdPartInitFailed   = 1010006
	UploadFailed          = 1010007
	ServerOperationFailed = 1010008
	PaymentRequestFailed  = 1010009
)
