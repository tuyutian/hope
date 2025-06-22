# golang logger

     Logger module based on uber zap package,
     which can be used to record operation logs in go applications
         1) Support automatic log cutting and maximum retention time
         2) Support log json formatting processing
         3) Support log output to file and terminal at the same time
         4) Support log printing level and log dyeing function
         
    基于uber zap封装而成的logger模块，可用于go应用中记录操作日志
        1）支持日志自动切割和最大保留时长
        2）支持日志json格式化处理
        3）支持日志同时输出到文件和终端
        4）支持日志打印级别和日志染色功能
        5）支持默认日志句柄调用，比如Info,Warn等方法调用

# how to use logger
```go
// 在go项目中引入logger包
import(
	"pkg/logger"
)

// 在main.go的main函数中初始化日志配置
// 日志输出采用zap框架实现日志json格式输出
// 日志level通过 logger.WithLogLevel(zap.DebugLevel) 设置
logger.Default(
    logger.WriteToFile(false), logger.WithStdout(true), // 将日志写到stdout
    logger.WithAddCaller(true), logger.WithLogLevel(zap.DebugLevel),
)

// 使用logger
logger.Info(context.Background(), "starting server",map[string]interface{}{
	"a":1,
	"b":2,
})

// 当然也支持zap格式
logger.Info(context.Background(), "starting server",
    "a",1,
    "b",2,
)

// 或者下面的格式
logger.Info(context.Background(), "starting server", zap.Int("pid", pid))
logger.Info(context.Background(), "starting server",
   zap.String("foo","abc"),
   zap.Int("b",2),
)
// 也支持下面的格式
logger.Info(context.Background(), "starting server", zap.Int("pid", 1),"b",2)

// 如果需要把所有的日志通过log_id串起来，可以通过下面的方式打印日志
import(
    "pkg/logger"
    "pkg/utils"
)

// ...省略日志初始化操作...
// 创建uuid作为日志log_id
logId := utils.Uuid()
ctx := context.Background()
ctx = context.WithValue(ctx, logger.XRequestID, logId)
// 下面的日志就会自动串联起来
logger.Info(ctx, "exec begin")
logger.Info(ctx, "foo")
```

日志格式如下：
```json
{"level":"info","time_local":"2025-03-24T10:04:36.962+0800","caller_line":"/web/go/hermes/cmd/web/main.go:34","msg":"starting server","pid":549,"local_time":"2025-03-24 10:04:36.962","hostname":"myapp.com","x-request-id":"2bd918fbf5d943758d3f34410eb4a8c0"}
```

# zap
    
    https://github.com/uber-go/zap
