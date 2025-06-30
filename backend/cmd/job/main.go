package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"backend/internal/application"
	"backend/internal/infras/config"
	"backend/internal/interfaces/job/handler"
	"backend/internal/interfaces/job/tasks"
	"backend/internal/providers"
	"backend/pkg/logger"
)

func main() {
	pid := os.Getpid()
	fmt.Printf("current asynq worker pid: %d\n", pid)

	// 初始化配置
	appConf := config.InitAppConfig()

	// 日志初始化
	logger.Default(
		logger.WriteToFile(true),
		logger.WithStdout(true),
		logger.WithAddCaller(true),
		logger.WithLogLevel(appConf.GetLogLevel()),
	)

	logger.Info(context.Background(), "starting asynq worker", zap.Int("pid", pid))

	// 初始化依赖
	db, err := config.NewDB("db_conf")
	if err != nil {
		log.Fatalf("db init error:%v", err)
	}

	redisClient, err := config.NewRedis("redis_conf")
	if err != nil {
		log.Fatalf("redis init error:%v", err)
	}

	server, err := config.NewAsynqServer("redis_conf")
	if err != nil {
		log.Fatalf("asynq server init error:%v", err)
	}

	// 初始化业务组件
	repos := providers.NewRepositories(db, redisClient, appConf)
	services := application.NewServices(repos)
	handlers := handler.InitHanders(services)

	// 注册任务处理器
	mux := asynq.NewServeMux()
	tasks.InitTask(mux, handlers)

	// 设置信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中启动 Asynq Server - 这样做完全正常！
	serverErr := make(chan error, 1)
	go func() {
		log.Println("🚀 Asynq worker started and ready to process jobs...")
		if err := server.Run(mux); err != nil {
			log.Printf("❌ Asynq server run error: %v", err)
			serverErr <- err
		}
	}()

	// 等待退出信号或服务错误
	select {
	case sig := <-ch:
		log.Printf("📨 Received signal: %s", sig.String())
	case err := <-serverErr:
		log.Printf("💥 Server error: %v", err)
		return
	}

	// 优雅关闭
	log.Println("🛑 Shutting down asynq worker...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Println("⏳ Waiting for current jobs to finish...")
		server.Shutdown() // 等待当前任务完成
		log.Println("✅ All jobs completed")
	}()

	select {
	case <-done:
		log.Println("✨ Asynq worker stopped gracefully")
	case <-ctx.Done():
		log.Println("⏰ Shutdown timeout, forcing exit")
	}

	log.Println("👋 Asynq worker exited")
}
