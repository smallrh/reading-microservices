@echo off
echo 🔧 Starting Development Environment...

REM 进入部署目录
cd deployments\docker

REM 启动基础设施服务
echo 🏗️ Starting infrastructure services...
docker compose up -d mysql redis consul

echo 🎉 Infrastructure is ready! You can now start your services manually:
echo.
echo Terminal 1 - User Service:
echo   cd user-service ^&^& go run main.go
echo.
echo Terminal 2 - Content Service:
echo   cd content-service ^&^& go run main.go
echo.
echo Terminal 3 - Reading Service:
echo   cd reading-service ^&^& go run main.go
echo.
echo Terminal 4 - Payment Service:
echo   cd payment-service ^&^& go run main.go
echo.
echo Terminal 5 - Notification Service:
echo   cd notification-service ^&^& go run main.go
echo.
echo Terminal 6 - Download Service:
echo   cd download-service ^&^& go run main.go
echo.
echo Terminal 7 - API Gateway:
echo   cd api-gateway ^&^& go run main.go
echo.
echo 📊 Management URLs:
echo   Consul UI: http://localhost:8500
echo   MySQL: localhost:3306 (root/password)
echo   Redis: localhost:6379

pause