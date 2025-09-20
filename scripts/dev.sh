#!/bin/bash

# 开发环境启动脚本

set -e

echo "🔧 Starting Development Environment..."

# 启动基础设施服务
echo "🏗️ Starting infrastructure services..."
cd "$(dirname "$0")/../deployments/docker"
docker-compose up -d mysql redis consul

# 等待基础设施服务启动
echo "⏳ Waiting for infrastructure services..."
sleep 15

# 检查基础服务健康状态
echo "🏥 Checking infrastructure health..."

# 检查MySQL
until docker exec reading-mysql mysqladmin ping -h"localhost" --silent; do
    echo "⏳ Waiting for MySQL to be ready..."
    sleep 2
done
echo "✅ MySQL is ready"

# 检查Redis
until docker exec reading-redis redis-cli ping | grep -q PONG; do
    echo "⏳ Waiting for Redis to be ready..."
    sleep 2
done
echo "✅ Redis is ready"

echo "🎉 Infrastructure is ready! You can now start your services manually:"
echo ""
echo "Terminal 1 - User Service:"
echo "  cd user-service && go run main.go"
echo ""
echo "Terminal 2 - Content Service:"
echo "  cd content-service && go run main.go"
echo ""
echo "Terminal 3 - Reading Service:"
echo "  cd reading-service && go run main.go"
echo ""
echo "Terminal 4 - Payment Service:"
echo "  cd payment-service && go run main.go"
echo ""
echo "Terminal 5 - Notification Service:"
echo "  cd notification-service && go run main.go"
echo ""
echo "Terminal 6 - Download Service:"
echo "  cd download-service && go run main.go"
echo ""
echo "Terminal 7 - API Gateway:"
echo "  cd api-gateway && go run main.go"
echo ""
echo "📊 Management URLs:"
echo "  Consul UI: http://localhost:8500"
echo "  MySQL: localhost:3306 (root/password)"
echo "  Redis: localhost:6379"