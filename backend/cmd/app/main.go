package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/infras/config"
	"backend/internal/interfaces/web/handler"
	routers "backend/internal/interfaces/web/router"
	"backend/internal/providers"
	"backend/pkg/logger"
)

func main() {
	pid := os.Getpid()
	fmt.Printf("current service pid: %d\n", pid)
	// 初始化配置
	appConf := config.InitAppConfig()
	// 日志输出采用zap框架实现日志json格式输出
	logger.Default(
		logger.WriteToFile(true), logger.WithStdout(true), // 将日志写到stdout
		logger.WithAddCaller(true), logger.WithLogLevel(appConf.GetLogLevel()),
	)

	logger.Info(context.Background(), "starting server", zap.Int("pid", pid))
	db, err := config.NewDB("db_conf")
	if err != nil {
		log.Fatalf("db init error:%v", err)
	}
	redisClient, err := config.NewRedis("redis_conf")
	if err != nil {
		log.Fatalf("redis init error:%v", err)
	}
	// 初始化repos
	repos := providers.NewRepositories(db, redisClient)
	// 初始化 handlers
	handlers := handler.InitHandlers(repos)
	// 初始化路由规则
	router := gin.New()
	// 注册路由规则
	routers.InitRouters(router, handlers)

}
