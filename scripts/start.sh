#!/bin/bash

# 读书应用微服务启动脚本

set -e

echo "🚀 Starting Reading Microservices..."

# 检查Docker是否运行
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# 进入部署目录
cd "$(dirname "$0")/../deployments/docker"

# 停止并移除现有容器
echo "🧹 Cleaning up existing containers..."
docker-compose down -v

# 构建并启动服务
echo "🔨 Building and starting services..."
docker-compose up -d --build

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

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"

    echo "Checking $name on port $port..."

    # 等待服务启动
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
        docker-compose logs "$name"
        exit 1
    fi
done

echo "🎉 All services are running successfully!"
echo ""
echo "📋 Service URLs:"
echo "  API Gateway: http://localhost:8080"
echo "  User Service: http://localhost:8081"
echo "  Content Service: http://localhost:8082"
echo "  Reading Service: http://localhost:8083"
echo "  Payment Service: http://localhost:8084"
echo "  Notification Service: http://localhost:8085"
echo "  Download Service: http://localhost:8086"
echo ""
echo "📊 Management URLs:"
echo "  Consul UI: http://localhost:8500"
echo ""
echo "🔍 To view logs: docker-compose logs -f [service-name]"
echo "🛑 To stop all services: docker-compose down"