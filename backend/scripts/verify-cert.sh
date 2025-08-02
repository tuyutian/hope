#!/bin/bash

echo "=== 证书挂载验证脚本 ==="
echo

echo "1. 检查主机证书文件..."
if [ -f "/etc/letsencrypt/live/api.protectifyapp.com/fullchain.pem" ]; then
    echo "✅ 主机证书文件存在"
    ls -la /etc/letsencrypt/live/api.protectifyapp.com/
else
    echo "❌ 主机证书文件不存在"
    echo "请先运行: sudo certbot certonly --standalone -d api.protectifyapp.com"
    exit 1
fi
echo

echo "2. 检查nginx容器内证书文件..."
docker exec -it hope-nginx ls -la /etc/letsencrypt/live/api.protectifyapp.com/
echo

echo "3. 测试nginx配置..."
docker exec -it hope-nginx nginx -t
echo

echo "4. 检查证书内容..."
echo "主机证书信息:"
openssl x509 -in /etc/letsencrypt/live/api.protectifyapp.com/fullchain.pem -text -noout | grep "Subject:"
echo

echo "5. 重新加载nginx配置..."
docker exec hope-nginx nginx -s reload
echo

echo "6. 测试HTTPS连接..."
curl -k -I https://api.protectifyapp.com/
echo

echo "验证完成！" 