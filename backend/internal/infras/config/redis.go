package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"backend/pkg/gredis"
)

// NewRedis 创建redis实例
func NewRedis(name string) (redis.UniversalClient, error) {
	redisConf := gredis.RedisConf{}
	err := conf.ReadSection(name, &redisConf)
	if err != nil {
		return nil, fmt.Errorf("failed to read config for %s section: %s", name, err)
	}

	client := redisConf.InitClient()
	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %s", err)
	}

	return client, nil
}
