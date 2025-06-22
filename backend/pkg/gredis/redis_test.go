package gredis

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	conf := RedisConf{
		Address:         "",
		Password:        "",
		DB:              0,
		MaxRetries:      2,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		PoolSize:        20,
		PoolTimeout:     5 * time.Second,
		MinIdleConns:    20,
		ConnMaxIdleTime: 1200 * time.Second,
		ConnMaxLifetime: 1800 * time.Second,
	}

	redisClient := conf.InitClient()
	key := "foo"
	err := redisClient.Set(context.Background(), key, "abc", 0).Err()
	if err != nil {
		log.Fatalf("redis set error:%v", err)
	}

	log.Println("redis set success")
}

func TestRedisCluster(t *testing.T) {
	// redis cluster test
	clusterConf := RedisClusterConf{
		AddressNodes: []string{
			"127.0.0.1:6391",
			"127.0.0.1:6392",
			"127.0.0.1:6393",
			"127.0.0.1:6394",
			"127.0.0.1:6395",
			"127.0.0.1:6396",
		},
		PoolSize:     10, // PoolSize applies per cluster node and not for the whole cluster.
		MaxRetries:   2,  // 重试次数
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second, // 底层默认3s
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
		MinIdleConns: 10,
	}

	cluster := clusterConf.InitClusterClient()
	defer cluster.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	str, err := cluster.Set(ctx, "username", "daheige", 1000*time.Second).Result()
	log.Println(str, err)

	str, err = cluster.Set(ctx, "myname", "daheige2", 1000*time.Second).Result()
	log.Println(str, err)

	log.Println(cluster.Get(ctx, "username").Result())
}
