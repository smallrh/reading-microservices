@echo off
echo 🚀 Starting Reading Microservices...

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not running. Please start Docker Desktop first.
    pause
    exit /b 1
)

REM 进入部署目录
cd deployments\docker

REM 启动所有服务
echo 🔨 Building and starting services...
docker compose up -d --build

if errorlevel 1 (
    echo ❌ Failed to start services. Please check Docker installation.
    pause
    exit /b 1
)

echo 🎉 All services are starting!
echo.
echo 📋 Service URLs:
echo   API Gateway: http://localhost:8080
echo   User Service: http://localhost:8081
echo   Content Service: http://localhost:8082
echo   Reading Service: http://localhost:8083
echo   Payment Service: http://localhost:8084
echo   Notification Service: http://localhost:8085
echo   Download Service: http://localhost:8086
echo   Dify AI Platform: http://localhost:8093
echo   Dify API: http://localhost:8091
echo.
echo 📊 Management URLs:
echo   Consul UI: http://localhost:8500
echo.
echo 🔍 To view logs: docker compose logs -f [service-name]
echo 🛑 To stop all services: docker compose down

pause