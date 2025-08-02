package apps

// AppData 存储应用程序的基本信息
type AppData struct {
	AppID     string `json:"app_id"`     // 应用ID
	AppKey    string `json:"app_key"`    // 应用Key
	AppSecret string `json:"app_secret"` // 应用Secret
}
