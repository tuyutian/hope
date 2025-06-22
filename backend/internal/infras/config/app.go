package config

import (
	"backend/pkg/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"path/filepath"
	"time"
)

// AppConfig 应用基本配置
type AppConfig struct {
	// 这里的 mapstructure tag标签用于指定配置文件的字段名字
	AppDebug bool   `mapstructure:"app_debug"` // 是否开启调试模式
	AppEnv   string `mapstructure:"app_env"`   // prod,test,local,dev

	MonitorPort  uint16        `mapstructure:"monitor_port"`  // metrics服务端口
	GrpcPort     uint16        `mapstructure:"grpc_port"`     // grpc 服务端口
	GracefulWait time.Duration `mapstructure:"graceful_wait"` // 平滑退出等待时间
	LogLevel     string        `mapstructure:"log_level"`     // 日志等级
}

// 配置文件读取的接口
var conf settings.Config

// InitAppConfig 初始化app config
// 这个函数的代码可以根据实际情况在main.go初始化
func InitAppConfig() *AppConfig {
	var err error
	// 读取配置文件，并初始化redis和mysql
	conf, err = loadConfig("./app.yaml")
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	appConfig := &AppConfig{}
	err = conf.ReadSection("app_conf", appConfig)
	if err != nil {
		log.Fatalln("failed to read config:", err)
	}

	// log.Println("read app_conf err: ", err)
	if appConfig.AppDebug {
		log.Println("app_conf:", appConfig)
	}

	return appConfig
}

// LoadConfig 加载配置文件
func loadConfig(path string) (settings.Config, error) {
	// 加载配置文件
	log.Println("config filename:", path, " dir:", filepath.Dir(path))
	c := settings.New(settings.WithConfigFile(path))
	err := c.Load()
	if err != nil {
		return nil, err
	}

	return c, nil
}
func (appConfig *AppConfig) GetLogLevel() zapcore.Level {
	switch appConfig.LogLevel {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	default:
		return zap.ErrorLevel
	}
}
