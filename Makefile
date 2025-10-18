.PHONY: help build up down logs ps clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all services
	docker compose build

up: ## Start all services
	docker compose up -d

up-build: ## Build and start all services
	docker compose up -d --build

down: ## Stop all services
	docker compose down

down-v: ## Stop all services and remove volumes
	docker compose down -v

logs: ## View logs from all services
	docker compose logs -f

logs-books: ## View logs from books service
	docker compose logs -f books-service

logs-users: ## View logs from users service
	docker compose logs -f users-service

logs-logging: ## View logs from logging service
	docker compose logs -f logging-service

ps: ## List running services
	docker compose ps

clean: ## Stop services and clean up
	docker compose down -v
	docker system prune -f

restart: ## Restart all services
	docker compose restart

test-books: ## Run tests for books service
	cd services/books-service && go test ./... -v -cover

test-users: ## Run tests for users service
	cd services/users-service && go test ./... -v -cover

test-logging: ## Run tests for logging service
	cd services/logging-service && go test ./... -v -cover

test-all: test-books test-users test-logging ## Run all tests

fmt-books: ## Format books service code
	cd services/books-service && gofmt -s -w .

fmt-users: ## Format users service code
	cd services/users-service && gofmt -s -w .

fmt-logging: ## Format logging service code
	cd services/logging-service && gofmt -s -w .

fmt-all: fmt-books fmt-users fmt-logging ## Format all code

health: ## Check health of all services
	@echo "Books Service:"
	@curl -s http://localhost:8081/health || echo "Not running"
	@echo "\nUsers Service:"
	@curl -s http://localhost:8082/health || echo "Not running"
	@echo "\nLogging Service:"
	@curl -s http://localhost:8084/health || echo "Not running"

ready: ## Check readiness of all services
	@echo "Books Service:"
	@curl -s http://localhost:8081/ready || echo "Not ready"
	@echo "\nUsers Service:"
	@curl -s http://localhost:8082/ready || echo "Not ready"
	@echo "\nLogging Service:"
	@curl -s http://localhost:8084/ready || echo "Not ready"
