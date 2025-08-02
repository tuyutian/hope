#!/bin/bash

echo "=== SSL证书设置脚本 (CentOS) ==="
echo

echo "1. 停止nginx容器..."
docker-compose stop nginx
echo

echo "2. 检查certbot是否安装..."
if ! command -v certbot &> /dev/null; then
    echo "安装certbot..."
    # CentOS 7
    if [ -f /etc/redhat-release ]; then
        version=$(cat /etc/redhat-release | grep -oE '[0-9]+\.[0-9]+' | head -1)
        if [[ $version == 7* ]]; then
            echo "检测到CentOS 7，安装EPEL和certbot..."
            sudo yum install -y epel-release
            sudo yum install -y certbot
        else
            echo "检测到CentOS 8+，使用dnf安装..."
            sudo dnf install -y certbot
        fi
    else
        echo "无法确定CentOS版本，尝试yum安装..."
        sudo yum install -y epel-release
        sudo yum install -y certbot
    fi
fi
echo

echo "3. 获取SSL证书..."
sudo certbot certonly --standalone -d api.protectifyapp.com
echo

echo "4. 检查证书文件..."
if [ -f "/etc/letsencrypt/live/api.protectifyapp.com/fullchain.pem" ]; then
    echo "证书获取成功！"
    ls -la /etc/letsencrypt/live/api.protectifyapp.com/
else
    echo "证书获取失败，请检查错误信息"
    exit 1
fi
echo

echo "5. 启动nginx容器..."
docker-compose start nginx
echo

echo "6. 检查nginx配置..."
docker exec -it hope-nginx nginx -t
echo

echo "7. 重新加载nginx配置..."
docker exec hope-nginx nginx -s reload
echo

echo "8. 测试HTTPS连接..."
curl -k https://api.protectifyapp.com/health
echo

echo "SSL证书设置完成！"
echo "现在可以通过 https://api.protectifyapp.com 访问您的API"
echo
echo "设置证书自动续期："
echo "sudo crontab -e"
echo "添加: 0 12 * * * /usr/bin/certbot renew --quiet && docker exec hope-nginx nginx -s reload" 