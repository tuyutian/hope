# Docker 中国镜像源配置指南

本项目已配置好 Docker 容器使用中国的镜像源，以提高在中国大陆地区的构建和拉取速度。

**注意**: 前端部署在 Cloudflare Pages 上，此 Docker 配置仅用于后端服务。

## 🏗️ 服务架构

本项目包含两个后端服务：

- **API 服务** (`backend-api`): 主要的 REST API 服务，处理前端请求
  - 端口: 8080 (API), 8090 (管理)
  - 容器名: `hope-go-api`
  
- **Job 服务** (`backend-job`): 后台任务处理服务，处理异步任务
  - 端口: 8091 (映射到容器内8090)
  - 容器名: `hope-go-job`

两个服务使用相同的 Docker 镜像，通过启动参数区分：
- API 服务运行: `./api`
- Job 服务运行: `./job`

## 🚀 快速开始

### 1. 自动配置（推荐）

运行配置脚本自动设置 Docker 镜像源：

```bash
# 设置执行权限
chmod +x docker-setup.sh

# 运行配置脚本
./docker-setup.sh
```

### 2. 手动配置

#### macOS Docker Desktop

创建或编辑 `~/.docker/daemon.json`：

```json
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://registry.docker-cn.com"
  ]
}
```

然后重启 Docker Desktop。

#### Linux 系统

创建或编辑 `/etc/docker/daemon.json`：

```json
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://registry.docker-cn.com"
  ]
}
```

重启 Docker 服务：

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

## 📦 项目构建

### 使用 Docker Compose（推荐）

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 单独构建

#### 后端 API 服务

```bash
cd backend
docker build -t hope-backend .
docker run -d -p 8080:8080 -p 8090:8090 --name hope-api hope-backend api
```

#### 后端 Job 服务

```bash
cd backend
docker build -t hope-backend .
docker run -d -p 8091:8090 --name hope-job hope-backend job
```

## 🔧 已配置的国内镜像源

### Docker Registry 镜像源

- **腾讯云**: `https://mirror.ccs.tencentyun.com`
- **中科大**: `https://docker.mirrors.ustc.edu.cn`
- **网易**: `https://hub-mirror.c.163.com`
- **Docker中国**: `https://registry.docker-cn.com`

### 基础镜像源

项目 Dockerfile 中使用的基础镜像均来自阿里云镜像仓库：

- Go: `registry.cn-hangzhou.aliyuncs.com/library/golang:1.21-alpine`
- Alpine: `registry.cn-hangzhou.aliyuncs.com/library/alpine:latest`
- Redis: `registry.cn-hangzhou.aliyuncs.com/library/redis:7-alpine`
- MySQL: `registry.cn-hangzhou.aliyuncs.com/library/mysql:8.0`

### 包管理器镜像源

#### Go 模块代理

```bash
GOPROXY=https://goproxy.cn,direct
```

#### Alpine 包管理器

```bash
sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
```

## 🛠 故障排除

### 验证配置是否生效

```bash
# 查看 Docker 镜像源配置
docker info | grep -A 10 "Registry Mirrors"

# 测试拉取镜像速度
time docker pull registry.cn-hangzhou.aliyuncs.com/library/alpine:latest
```

### 常见问题

1. **镜像拉取仍然很慢**
   - 检查网络连接
   - 尝试不同的镜像源
   - 确认 Docker 配置是否生效

2. **构建失败**
   - 检查 Dockerfile 语法
   - 确认网络访问权限
   - 查看构建日志定位问题

3. **权限问题**
   - 确保 Docker 服务正在运行
   - 检查用户是否在 docker 组中

## 📚 相关链接

- [Docker 官方文档](https://docs.docker.com/)
- [阿里云容器镜像服务](https://cr.console.aliyun.com/)
- [腾讯云容器镜像服务](https://console.cloud.tencent.com/tcr)
- [中科大 Docker 镜像源](https://mirrors.ustc.edu.cn/help/dockerhub.html)

## 🤝 贡献

如果您发现更好的镜像源或有改进建议，欢迎提交 Issue 或 Pull Request。