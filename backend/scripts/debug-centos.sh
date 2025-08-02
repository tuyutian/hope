#!/bin/bash

echo "=== CentOS环境诊断脚本 ==="
echo

echo "1. 检查系统信息..."
cat /etc/redhat-release
echo

echo "2. 检查防火墙状态..."
sudo systemctl status firewalld
echo

echo "3. 检查SELinux状态..."
getenforce
echo

echo "4. 检查端口监听..."
netstat -tlnp | grep :80
netstat -tlnp | grep :443
echo

echo "5. 检查容器状态..."
docker-compose ps
echo

echo "6. 检查nginx容器日志..."
docker logs --tail 10 hope-nginx
echo

echo "7. 检查certbot状态..."
if command -v certbot &> /dev/null; then
    echo "certbot已安装"
    certbot --version
else
    echo "certbot未安装"
fi
echo

echo "8. 检查证书文件..."
if [ -d "/etc/letsencrypt/live/api.protectifyapp.com" ]; then
    echo "证书目录存在"
    ls -la /etc/letsencrypt/live/api.protectifyapp.com/
else
    echo "证书目录不存在"
fi
echo

echo "9. 测试HTTP访问..."
curl -I http://api.protectifyapp.com/
echo

echo "10. 测试HTTPS访问..."
curl -k -I https://api.protectifyapp.com/
echo

echo "诊断完成！" 