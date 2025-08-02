# 禁用 CGO 以避免 GLIBC 版本依赖
.ONESHELL:
export CGO_ENABLED=0 GOOS=linux
# 定义默认 Go 命令
GO_CMD ?= go

# 检查是否传入了 -v 参数来指定版本
ifdef GOVERSION
	GO_CMD = go$(GOVERSION)
endif
# docker容器相关配置
# 镜像版本号
VERSION :=v1.0
# 镜像名字
IMAGE_NAME :=hope-go
# 容器运行名字
API_NAME :=hope-go-api
JOB_NAME :=hope-go-job

# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
#rm -rf $(CURDIR)/cmd $(CURDIR)/pkg $(CURDIR)/internal $(CURDIR)/server $(CURDIR)/store $(CURDIR)/vendor $(CURDIR)/test
# 检查指定的 Go 版本是否存在
check-go-version:
ifdef GOVERSION
	@command -v go$(GOVERSION) > /dev/null 2>&1 || { echo "go$(GOVERSION) 未找到，请先安装"; exit 1; }
else
	@command -v go > /dev/null 2>&1 || { echo "go 未找到，请先安装"; exit 1; }
endif

# 将 check-go-version 添加为依赖
test: check-go-version
	$(GO_CMD) version


all:setup build-api build-job clean

setup:
	@echo "install dependency"
	@$(GO_CMD) mod tidy
build-api:
	CGO_ENABLED=0 $(GO_CMD) build -ldflags="-w -s -extldflags '-static'" -o admin $(CURDIR)/cmd/app/main.go
	@echo "build api success"
build-job:
	CGO_ENABLED=0 $(GO_CMD) build -ldflags="-w -s -extldflags '-static'" -o webhook $(CURDIR)/cmd/job/main.go
	@echo "build job success"
clean:
	@rm -rf $(CURDIR)/internal $(CURDIR)/pkg $(CURDIR)/server $(CURDIR)/store $(CURDIR)/test $(CURDIR)/tmp $(CURDIR)/.git
	@echo "clear success"

help:
	@echo "Hope项目 Docker 构建帮助"
	@echo "=========================="
	@echo ""
	@echo ""
	@echo "传统构建命令:"
	@echo "  make build       - 构建 Docker 镜像"
	@echo "  make run         - 运行容器"
	@echo "  make stop        - 停止容器"
	@echo "  make remove      - 删除容器"
	@echo ""
	@echo "Docker Compose 命令:"
	@echo "  make compose-build  - 构建所有服务"
	@echo "  make compose-up     - 启动所有服务"
	@echo "  make compose-down   - 停止所有服务"
	@echo "  make compose-logs   - 查看服务日志"
	@echo ""
	@echo "快速开始:"
	@echo "  1. 配置 Docker 镜像源: ../docker-setup.sh"
	@echo "  2. 构建并启动: make compose-up"

# 构建容器镜像
build:
	@echo "docker build start for "${IMAGE_NAME}:${VERSION}
	@docker build --build-arg  -t ${IMAGE_NAME}:${VERSION} . -f Dockerfile

compose-up:
	@echo "启动所有服务..."
	@cd .. && docker-compose up -d

compose-down:
	@echo "停止所有服务..."
	@cd .. && docker-compose down

compose-logs:
	@echo "查看服务日志..."
	@cd .. && docker-compose logs -f

# 运行容器
run: run_api run_job
run_api:
	@docker run -itd \
       --name ${API_NAME} \
       -p 8080:8080 -p 8090:8090 \
       ${API_NAME}:${VERSION}
	@echo "${API_NAME} run success"
run_job:
	@docker run -itd \
       --name ${JOB_NAME} \
       -p 8091:8090 \
       ${IMAGE_NAME}:${VERSION}

	@echo "${JOB_NAME} run success"

# 重新运行
rerun: remove run_api run_job

# 重新构建和运行
rebuild-run: build rerun

# 停止容器
stop:
	docker stop ${API_NAME} ${JOB_NAME}

# 容器重启
restart: restart-api restart-job

restart-api:
	@echo "${API_NAME} restart success"
	docker restart ${API_NAME}
restart-job:
	@echo "${JOB_NAME} restart success"
	docker restart  ${JOB_NAME}

# 移除容器
remove:
	docker rm -f ${API_NAME} ${JOB_NAME}
# 进入容器
exec-api:
	docker exec -it ${API_NAME} /bin/bash
exec-job:
	docker exec -it ${JOB_NAME} /bin/bash

# 查看容器运行日志
logs:
	docker logs -f ${API_NAME}

.PHONY: build-job build-api clean setup all
