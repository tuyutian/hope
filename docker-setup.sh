#!/bin/bash

# Docker 中国镜像源配置脚本
# 支持 macOS 和 Linux

echo "🚀 开始配置 Docker 中国镜像源..."

# 检测操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)     MACHINE=Linux;;
    Darwin*)    MACHINE=Mac;;
    *)          MACHINE="UNKNOWN:${OS}"
esac

echo "检测到操作系统: ${MACHINE}"

# 配置 Docker Registry 镜像源
configure_docker_registry() {
    echo "📦 配置 Docker Registry 镜像源..."
    
    if [[ "$MACHINE" == "Mac" ]]; then
        # macOS Docker Desktop 配置
        DOCKER_CONFIG_DIR="$HOME/.docker"
        mkdir -p "$DOCKER_CONFIG_DIR"
        
        cat > "$DOCKER_CONFIG_DIR/daemon.json" << EOF
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://registry.docker-cn.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
        echo "✅ macOS Docker Desktop 镜像源配置完成"
        echo "请重启 Docker Desktop 使配置生效"
        
    elif [[ "$MACHINE" == "Linux" ]]; then
        # Linux 系统配置
        sudo mkdir -p /etc/docker
        
        sudo tee /etc/docker/daemon.json > /dev/null << EOF
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://registry.docker-cn.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
        
        echo "✅ Linux Docker 镜像源配置完成"
        echo "正在重启 Docker 服务..."
        sudo systemctl daemon-reload
        sudo systemctl restart docker
        echo "✅ Docker 服务重启完成"
    fi
}

# 配置构建时的包管理器镜像源
configure_build_sources() {
    echo "🔧 配置构建时包管理器镜像源..."
    
    # 检查是否存在 backend/Dockerfile
    if [[ -f "backend/Dockerfile" ]]; then
        echo "✅ 后端 Dockerfile 已存在，包含 Go 模块代理配置"
    else
        echo "⚠️  后端 Dockerfile 不存在"
    fi
    
    # 检查是否存在 frontend/Dockerfile
    if [[ -f "frontend/Dockerfile" ]]; then
        echo "✅ 前端 Dockerfile 已存在，包含 npm 镜像源配置"
    else
        echo "⚠️  前端 Dockerfile 不存在"
    fi
}

# 验证配置
verify_configuration() {
    echo "🔍 验证 Docker 配置..."
    
    # 测试拉取镜像
    echo "测试拉取阿里云镜像..."
    docker pull registry.cn-hangzhou.aliyuncs.com/library/alpine:latest
    
    if [[ $? -eq 0 ]]; then
        echo "✅ 镜像拉取测试成功"
        docker rmi registry.cn-hangzhou.aliyuncs.com/library/alpine:latest
    else
        echo "❌ 镜像拉取测试失败"
    fi
    
    # 显示当前配置
    echo "📋 当前 Docker 配置信息:"
    docker info | grep -A 10 "Registry Mirrors" || echo "未找到镜像源配置信息"
}

# 提供使用说明
show_usage() {
    echo ""
    echo "🎯 使用说明:"
    echo "1. 构建后端镜像: cd backend && docker build -t hope-backend ."
    echo "2. 构建前端镜像: cd frontend && docker build -t hope-frontend ."
    echo "3. 使用 docker-compose: docker-compose up -d"
    echo ""
    echo "📚 常用的中国镜像源:"
    echo "- 腾讯云: https://mirror.ccs.tencentyun.com"
    echo "- 中科大: https://docker.mirrors.ustc.edu.cn"
    echo "- 网易: https://hub-mirror.c.163.com"
    echo "- 阿里云: https://registry.cn-hangzhou.aliyuncs.com"
    echo ""
    echo "🛠  如需手动配置，请编辑:"
    if [[ "$MACHINE" == "Mac" ]]; then
        echo "- macOS: ~/.docker/daemon.json"
    else
        echo "- Linux: /etc/docker/daemon.json"
    fi
}

# 主执行流程
main() {
    configure_docker_registry
    configure_build_sources
    verify_configuration
    show_usage
    
    echo ""
    echo "🎉 Docker 中国镜像源配置完成!"
}

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 执行主函数
main