# 使用阿里云镜像的 Go 基础镜像
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置 Go 模块代理和环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# 设置 Alpine 包管理器国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git ca-certificates tzdata

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 构建两个应用
RUN go build -ldflags="-w -s -extldflags '-static'" -o api ./cmd/app/main.go && \
    go build -ldflags="-w -s -extldflags '-static'" -o job ./cmd/job/main.go

# 使用阿里云镜像的最小基础镜像
FROM alpine:latest

# 设置 Alpine 包管理器国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata

# 创建工作目录和日志目录
WORKDIR /app
RUN mkdir -p /app/logs

# 从builder阶段复制可执行文件
COPY --from=builder /app/api /app/api
COPY --from=builder /app/job /app/job

# 复制配置文件（如果存在）
COPY --from=builder /app/app.example.yaml /app/app.yaml

# 创建启动脚本
COPY ./scripts/entrypoint.sh /app/start.sh

# 设置文件权限
RUN chmod +x /app/start.sh

# 暴露端口
EXPOSE 8080 8090

# 启动脚本，不设置默认参数，让docker-compose的command能够正确传递
ENTRYPOINT ["./start.sh"]
