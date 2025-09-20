@echo off
echo 🔧 Setting environment variables for local development...

set REDIS_HOST=localhost
set REDIS_PORT=6379
set REDIS_DB=2
set SERVER_HOST=localhost
set SERVER_PORT=8080
set JWT_SECRET=reading-app-secret-key
set CONSUL_HOST=localhost
set CONSUL_PORT=8500

echo 🚀 Starting API Gateway locally...
cd api-gateway
go run main.go