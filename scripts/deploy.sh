#!/bin/bash

# 生产环境部署脚本 - Linux

set -e

echo "🚀 Deploying Reading Microservices to Production..."

# 检查是否为root用户或有sudo权限
if [[ $EUID -eq 0 ]]; then
   echo "⚠️  Running as root user"
elif ! sudo -n true 2>/dev/null; then
   echo "❌ This script requires sudo privileges"
   exit 1
fi

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "✅ Docker installed. Please logout and login again, then rerun this script."
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Installing..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "✅ Docker Compose installed"
fi

# 检查Docker服务状态
if ! sudo systemctl is-active --quiet docker; then
    echo "🔧 Starting Docker service..."
    sudo systemctl start docker
    sudo systemctl enable docker
fi

# 进入部署目录
cd "$(dirname "$0")/../deployments/docker"

# 停止现有服务
echo "🛑 Stopping existing services..."
sudo docker-compose down -v || true

# 清理旧镜像
echo "🧹 Cleaning up old images..."
sudo docker system prune -af || true

# 创建必要的目录
echo "📁 Creating necessary directories..."
sudo mkdir -p /var/lib/reading-app/{mysql,redis,downloads,logs}
sudo chown -R $(whoami):$(whoami) /var/lib/reading-app

# 设置环境变量（生产环境）
export ENVIRONMENT=production
export DATABASE_HOST=mysql
export DATABASE_PASSWORD=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 64)

# 生成生产环境配置
cat > .env.production << EOF
# Production Environment Variables
ENVIRONMENT=production

# Database
DATABASE_HOST=mysql
DATABASE_PORT=3306
DATABASE_USERNAME=root
DATABASE_PASSWORD=${DATABASE_PASSWORD}
DATABASE_DATABASE=reading_app

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRES_IN=86400

# Consul
CONSUL_HOST=consul
CONSUL_PORT=8500

# Log Level
LOG_LEVEL=info
EOF

echo "🔐 Generated production secrets in .env.production"
echo "⚠️  Please save these credentials securely!"

# 构建并启动服务
echo "🔨 Building and starting production services..."
sudo docker-compose --env-file .env.production up -d --build

# 等待服务启动
echo "⏳ Waiting for services to start..."
sleep 60

# 检查服务健康状态
echo "🏥 Checking service health..."

services=(
    "api-gateway:8080"
    "user-service:8081"
    "content-service:8082"
    "reading-service:8083"
    "payment-service:8084"
    "notification-service:8085"
    "download-service:8086"
)

all_healthy=true

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"

    echo "Checking $name on port $port..."

    timeout=120
    while [ $timeout -gt 0 ]; do
        if curl -f -s "http://localhost:$port/health" > /dev/null; then
            echo "✅ $name is healthy"
            break
        fi
        echo "⏳ Waiting for $name to be ready..."
        sleep 3
        timeout=$((timeout-3))
    done

    if [ $timeout -le 0 ]; then
        echo "❌ $name failed to start properly"
        echo "📋 Service logs:"
        sudo docker-compose logs "$name"
        all_healthy=false
    fi
done

if [ "$all_healthy" = true ]; then
    echo "🎉 All services deployed successfully!"
    echo ""
    echo "📋 Production Service URLs:"
    echo "  API Gateway: http://$(hostname -I | awk '{print $1}'):8080"
    echo "  Consul UI: http://$(hostname -I | awk '{print $1}'):8500"
    echo ""
    echo "🔍 To view logs: sudo docker-compose logs -f [service-name]"
    echo "🛑 To stop all services: sudo docker-compose down"
    echo "🔄 To restart services: sudo docker-compose restart [service-name]"
    echo ""
    echo "📁 Data directories:"
    echo "  MySQL: /var/lib/reading-app/mysql"
    echo "  Redis: /var/lib/reading-app/redis"
    echo "  Downloads: /var/lib/reading-app/downloads"
    echo "  Logs: /var/lib/reading-app/logs"
else
    echo "❌ Some services failed to start. Please check the logs."
    exit 1
fi