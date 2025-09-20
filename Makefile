.PHONY: help build start stop dev clean test

help: ## Show help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	@echo "🔨 Building all services..."
	@cd api-gateway && go mod tidy && go build -o bin/gateway main.go
	@cd user-service && go mod tidy && go build -o bin/user-service main.go
	@cd content-service && go mod tidy && go build -o bin/content-service main.go
	@cd reading-service && go mod tidy && go build -o bin/reading-service main.go
	@cd payment-service && go mod tidy && go build -o bin/payment-service main.go
	@cd notification-service && go mod tidy && go build -o bin/notification-service main.go
	@cd download-service && go mod tidy && go build -o bin/download-service main.go
	@echo "✅ All services built successfully"

start: ## Start all services using Docker Compose
ifeq ($(OS),Windows_NT)
	@cd deployments\docker && docker compose up -d --build
else
	@cd deployments/docker && docker-compose up -d --build
endif

stop: ## Stop all services
ifeq ($(OS),Windows_NT)
	@cd deployments\docker && docker compose down
else
	@cd deployments/docker && docker-compose down
endif

dev: ## Start development environment (infrastructure only)
ifeq ($(OS),Windows_NT)
	@cd deployments\docker && docker compose up -d mysql redis consul
else
	@cd deployments/docker && docker-compose up -d mysql redis consul
endif

clean: ## Clean up Docker resources
	@echo "🧹 Cleaning up Docker resources..."
	@cd deployments/docker && docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "✅ Cleanup completed"

test: ## Run tests for all services
	@echo "🧪 Running tests..."
	@cd shared && go test ./...
	@cd api-gateway && go test ./...
	@cd user-service && go test ./...
	@cd content-service && go test ./...
	@cd reading-service && go test ./...
	@cd payment-service && go test ./...
	@cd notification-service && go test ./...
	@cd download-service && go test ./...
	@echo "✅ All tests passed"

logs: ## Show logs from all services
	@cd deployments/docker && docker-compose logs -f

status: ## Show status of all services
	@echo "📊 Service Status:"
	@curl -s http://localhost:8080/status 2>/dev/null | jq . || echo "API Gateway not accessible"

init: ## Initialize Go modules for all services
	@echo "📦 Initializing Go modules..."
	@cd shared && go mod init reading-microservices/shared && go mod tidy
	@cd api-gateway && go mod tidy
	@cd user-service && go mod tidy
	@cd content-service && go mod tidy
	@cd reading-service && go mod tidy
	@cd payment-service && go mod tidy
	@cd notification-service && go mod tidy
	@cd download-service && go mod tidy
	@echo "✅ Go modules initialized"