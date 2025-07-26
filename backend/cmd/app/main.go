package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/application"
	"backend/internal/infras/config"
	"backend/internal/interfaces/web/handler"
	"backend/internal/interfaces/web/middleware"
	"backend/internal/interfaces/web/routers"
	"backend/internal/providers"
	"backend/pkg/logger"
	"backend/pkg/monitor"
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
	ossClient, bucketName, err := config.NewAliyunOSS("aliyun_oss")
	if err != nil {
		log.Fatalf("oss init error:%v", err)
	}
	asynqClient, err := config.NewAsynqClient("redis_conf")
	// 初始化repos
	repos := providers.NewRepositories(db, redisClient, appConf, providers.WithOssRepo(ossClient, bucketName), providers.WithAsynqRepo(asynqClient))
	// 初始化服务
	services := application.NewServices(repos)
	// 初始化 handlers
	handlers := handler.InitHandlers(services, repos)
	// 初始化 middlewares
	// init middleware and routers
	middlewares := &routers.Middleware{
		RequestWare:        &middleware.RequestWare{},
		CorsWare:           &middleware.CorsWare{},
		AuthWare:           middleware.NewAuthWare(services.UserService, services.AppService, repos),
		ShopifyGraphqlWare: middleware.NewShopifyGraphqlWare(repos, services.UserService),
	}
	// 初始化路由规则
	router := gin.New()
	// 注册路由规则
	routers.InitRouters(router, handlers, middlewares)
	router.GET("/test/token", func(c *gin.Context) {
		token := services.UserService.GenerateTestToken(c.Request.Context(), 1)
		c.JSON(200, gin.H{
			"message": "hello world",
			"token":   token,
		})
	})
	router.GET("/test/sync_user/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		initUserTask, _ := repos.AsyncRepo.InitUserTask(c.Request.Context(), int64(id))
		c.JSON(200, gin.H{
			"message": "hello sync",
			"task":    initUserTask,
		})
	})
	// http server设置
	server := &http.Server{
		Handler:           router,
		Addr:              fmt.Sprintf("0.0.0.0:%d", appConf.AppPort),
		IdleTimeout:       20 * time.Second, // tcp idle time
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	// 在独立携程中运行
	log.Println("server run on: ", appConf.AppPort)
	go func() {
		defer logger.Recover(context.Background(), "server start panic")
		if err2 := server.ListenAndServe(); err2 != nil {
			if !errors.Is(err2, http.ErrServerClosed) {
				log.Println(context.Background(), "server close error", map[string]interface{}{
					"trace_error": err2.Error(),
				})

				log.Println("server close error:", err2)
				return
			}

			log.Println("server will exit...")
		}
	}()

	// 初始化prometheus和pprof
	// 访问地址：http://localhost:8090/metrics
	// 访问地址：http://localhost:8090/debug/pprof/
	monitor.InitMonitor(appConf.MonitorPort, true)

	// server平滑重启
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	sig := <-ch
	log.Println("exit signal: ", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), appConf.GracefulWait)
	defer cancel()

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		if err2 := server.Shutdown(ctx); err2 != nil {
			log.Println("server shutdown error", map[string]interface{}{
				"trace_error": err2.Error(),
			})
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("context deadline timeout")
	case <-done:
	}

	log.Println("server shutting down")
}
