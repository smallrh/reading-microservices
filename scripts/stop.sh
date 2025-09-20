#!/bin/bash

# 停止所有服务的脚本

set -e

echo "🛑 Stopping Reading Microservices..."

# 进入部署目录
cd "$(dirname "$0")/../deployments/docker"

# 停止并移除容器
docker-compose down

echo "✅ All services have been stopped."

# 可选：清理未使用的Docker资源
read -p "Do you want to clean up unused Docker resources? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🧹 Cleaning up Docker resources..."
    docker system prune -f
    echo "✅ Docker resources cleaned up."
fi