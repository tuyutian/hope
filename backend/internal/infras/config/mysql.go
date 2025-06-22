package config

import (
	"fmt"

	"xorm.io/xorm"

	"backend/pkg/gxorm"
)

// NewDB 根据配置文件配置的名字获取DB句柄
func NewDB(name string) (*xorm.Engine, error) {
	dbConfig := gxorm.DbConf{}
	err := conf.ReadSection(name, &dbConfig)
	// log.Printf("db conf:v%\n", dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to read config for %s section: %s", name, err)
	}

	// 建立mysql连接
	db, err := dbConfig.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection for %s section: %s", name, err)
	}

	return db, nil
}
