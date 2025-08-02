# Nginx 配置说明

## 目录结构
```
nginx/
├── nginx.conf                    # 主配置文件
├── conf.d/                       # 站点配置文件目录
│   └── api.protectifyapp.com.conf # API域名配置
├── ssl/                         # SSL证书目录（可选）
└── README.md                    # 说明文档
```

## 配置说明

### 1. HTTP 配置
- 监听80端口
- 域名：api.protectifyapp.com
- 反向代理到后端API服务（backend-api:8080）

### 2. 安全特性
- 添加了安全响应头
- 配置了Gzip压缩
- 设置了适当的缓存策略

### 3. 健康检查
- 访问 `/health` 端点可以检查服务状态

## 使用方法

### 启动服务
```bash
# 启动所有服务（包括nginx）
docker-compose up -d

# 仅启动nginx
docker-compose up -d nginx
```

### 检查配置
```bash
# 进入nginx容器检查配置
docker exec -it hope-nginx nginx -t

# 查看nginx日志
docker logs hope-nginx
```

### 重新加载配置
```bash
# 重新加载nginx配置（不重启容器）
docker exec hope-nginx nginx -s reload
```

## HTTPS 配置（可选）

如果需要启用HTTPS：

1. 将SSL证书文件放入 `nginx/ssl/` 目录
2. 取消注释 `api.protectifyapp.com.conf` 中的HTTPS配置
3. 修改证书文件路径
4. 重启nginx服务

## 故障排除

### 1. 无法访问
- 检查DNS解析是否正确
- 确认阿里云安全组开放了80端口
- 查看nginx日志：`docker logs hope-nginx`

### 2. 502错误
- 检查后端API服务是否正常运行
- 确认容器网络连接正常
- 查看后端服务日志：`docker logs hope-go-api`

### 3. 配置修改后不生效
- 重新加载nginx配置：`docker exec hope-nginx nginx -s reload`
- 或重启nginx容器：`docker-compose restart nginx` 