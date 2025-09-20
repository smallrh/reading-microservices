@echo off
echo 🛑 Stopping Reading Microservices...

REM 进入部署目录
cd deployments\docker

REM 停止所有服务
docker compose down

echo ✅ All services have been stopped.

pause