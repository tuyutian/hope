package logger

import (
	"context"
	"log"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"

	"backend/pkg/ctxkeys"
	"backend/pkg/utils"
)

// TestLogger test logger.
func TestLogger(t *testing.T) {
	// 对于option 下面的可以根据实际情况使用
	var logger = New(
		WithLogDir("./logs"),
		WithLogFilename("zap.log"),
		WithStdout(true), // 生产环境，建议尽量输出到stdout，避免输出到文件中
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		WithHostname("myapp.com"),
	)

	reqId := utils.Uuid()
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.XRequestID, reqId)
	logger.Info(ctx, "hello", map[string]interface{}{
		"a": 1,
		"b": 12,
	})

	logger.Error(ctx, "exec error", zap.Any("details", map[string]interface{}{
		"name": "zap",
		"age":  30,
	}))

	logger.Debug(ctx, "test abc", nil)

	logger.Warn(ctx, "run warning", "key", 12)
	logger.DPanic(ctx, "exec panic but not exit", "stack", string(debug.Stack()))

	logger.Info(ctx, "abc")

	go func() {
		defer logger.Recover(ctx, "exec panic", "key", 123)

		x := 1
		log.Println("x = ", x)
		// panic(1111)
		logger.Panic(ctx, "current goroutine exit")

	}()

	time.Sleep(3 * time.Second)
	log.Println("exit...")
}

// TestNewLogSugar test log sugar.
func TestNewLogSugar(t *testing.T) {
	// 测试log sugar方法
	logSugar := NewLogSugar(WithLogDir("./logs"),
		WithLogFilename("zap-sugar.log"),
		WithStdout(true), // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		WithHostname("myapp.com"),
	)

	logSugar.Info("abc", 123, "info", "sugar hello")
	logSugar.Error("a", 234, "x", "sugar hello world")
}

/*
*
BenchmarkNew 批量测试日志写入
BenchmarkNew-12    	   18445	     72602 ns/op
*/
func BenchmarkNew(b *testing.B) {
	// 对于option 下面的可以根据实际情况使用
	var logger = New(
		WithLogDir("./logs"),
		WithLogFilename("zap-bench.log"),
		WithStdout(true), // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		// WithHostname("myapp.com"),
	)

	reqId := utils.Uuid()
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.XRequestID, reqId)
	logger.Info(ctx, "exec begin")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "hello", "index", i)
		logger.Error(ctx, "exec error", "abc", 1, "e", "zap is fast")
		logger.Info(ctx, "exec map", map[string]interface{}{
			"a": 1,
			"b": 123.23,
			"c": "hello,go",
			"e": []string{"f", "g", "higk"},
			"f": []int{1, 2, 3, i},
		})
	}

	logger.Info(ctx, "exec end", "cost_time", time.Since(start).Seconds())
}
