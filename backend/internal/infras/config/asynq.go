package config

import (
	"fmt"

	"github.com/hibiken/asynq"

	"backend/pkg/gredis"
)

const (
	SendProduct         = "task:send_product"
	SendInitUser        = "task:send_init_user"
	SendOrder           = "task:send_order"
	SendUpdateProduct   = "task:send_update_product"
	SendOrderStatistics = "task:send_order_statistics"
	SendDelProduct      = "task:send_delete_product"
)

func NewAsynqServer(name string) (*asynq.Server, error) {
	redisConf := gredis.RedisConf{}
	err := conf.ReadSection(name, &redisConf)
	if err != nil {
		return nil, fmt.Errorf("failed to read config for %s section: %s", name, err)
	}

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisConf.Address,
			Password: redisConf.Password, // no password set
			DB:       1,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)
	return server, nil
}
func NewAsynqClient(name string) (*asynq.Client, error) {
	redisConf := gredis.RedisConf{}
	err := conf.ReadSection(name, &redisConf)
	if err != nil {
		return nil, fmt.Errorf("failed to read config for %s section: %s", name, err)
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisConf.Address,
		Password: redisConf.Password, // no password set
		DB:       1,
	})

	return client, nil
}
