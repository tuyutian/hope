package settings

import (
	"log"
	"path/filepath"
	"testing"
)

// 假设app_test.yaml配置如下
// app_conf:
//    app_debug: true
//    app_env: local

func TestSettingsRead(t *testing.T) {
	// AppConfig 应用基本配置
	type AppConfig struct {
		// 这里的 mapstructure tag标签用于指定配置文件的字段名字
		AppDebug bool   `mapstructure:"app_debug"` // 是否开启调试模式
		AppEnv   string `mapstructure:"app_env"`   //  prod,test,local,dev
	}

	path := "./app_test.yaml"
	log.Println("config filename:", path, " dir:", filepath.Dir(path))
	c := New(WithConfigFile(path))
	err := c.Load()
	if err != nil {
		log.Fatalf("load config file error:%v", err)
	}

	var appConfig AppConfig
	err = c.ReadSection("app_conf", &appConfig)
	if err != nil {
		log.Fatalf("read app_conf error:%v", err)
	}

	log.Printf("app_conf:%+v", appConfig)
}

/*
2025/03/21 11:02:02 config filename: ./app_test.yaml  dir: .
2025/03/21 11:02:02 app_conf:{AppDebug:true AppEnv:local}
--- PASS: TestSettingsRead (0.00s)
PASS
*/
