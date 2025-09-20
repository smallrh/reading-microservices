#!/bin/bash

# 检查服务状态脚本

echo "📊 Reading Microservices Status Check"
echo "====================================="

# 检查Docker服务
if ! systemctl is-active --quiet docker; then
    echo "❌ Docker service is not running"
    exit 1
else
    echo "✅ Docker service is running"
fi

cd "$(dirname "$0")/../deployments/docker"

# 检查容器状态
echo ""
echo "📦 Container Status:"
docker-compose ps

echo ""
echo "🏥 Health Check Results:"

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

    if curl -f -s "http://localhost:$port/health" > /dev/null; then
        echo "✅ $name (port $port) - Healthy"
    else
        echo "❌ $name (port $port) - Unhealthy or not responding"
    fi
done

echo ""
echo "📈 Resource Usage:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

echo ""
echo "🔗 Service URLs:"
echo "  API Gateway: http://localhost:8080"
echo "  Consul UI: http://localhost:8500"