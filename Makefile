# ==========================================
# CS2 P2P Skins Trading Platform - Makefile
# ==========================================

.PHONY: help setup up down logs build rebuild clean test ps health backup restore

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)CS2 P2P Skins Trading Platform - Docker Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

setup: ## Initial setup - copy .env file
	@echo "$(BLUE)Setting up environment...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.docker.example .env; \
		echo "$(GREEN)✓$(NC) Created .env file from .env.docker.example"; \
		echo "$(RED)⚠$(NC)  Please edit .env and add your STEAM_API_KEY"; \
	else \
		echo "$(GREEN)✓$(NC) .env file already exists"; \
	fi

up: ## Start all services
	@echo "$(BLUE)Starting all services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓$(NC) All services started!"
	@echo ""
	@echo "Access the application:"
	@echo "  Frontend:  http://localhost:3000"
	@echo "  Backend:   http://localhost:8080"
	@echo "  MinIO:     http://localhost:9001"

down: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓$(NC) All services stopped"

logs: ## Show logs from all services
	@docker-compose logs -f

logs-backend: ## Show backend logs
	@docker-compose logs -f backend

logs-frontend: ## Show frontend logs
	@docker-compose logs -f frontend

build: ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	@docker-compose build
	@echo "$(GREEN)✓$(NC) Build complete"

rebuild: ## Rebuild and restart all services
	@echo "$(BLUE)Rebuilding and restarting...$(NC)"
	@docker-compose up -d --build --force-recreate
	@echo "$(GREEN)✓$(NC) Services rebuilt and restarted"

clean: ## Stop and remove all containers, networks, and volumes
	@echo "$(RED)WARNING: This will delete all data!$(NC)"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	@docker-compose down -v --rmi local
	@echo "$(GREEN)✓$(NC) All containers, volumes, and images removed"

test: ## Run tests (if available)
	@echo "$(BLUE)Running tests...$(NC)"
	@docker-compose exec backend go test ./...
	@echo "$(GREEN)✓$(NC) Tests complete"

ps: ## Show status of all services
	@docker-compose ps

health: ## Check health status of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@docker inspect --format='{{.Name}}: {{.State.Health.Status}}' $$(docker-compose ps -q) 2>/dev/null || echo "Health checks not available"

dev: setup ## Setup and start in development mode
	@echo "$(BLUE)Starting in development mode...$(NC)"
	@docker-compose up -d --build
	@echo "$(GREEN)✓$(NC) Development environment ready!"
	@make logs

restart: ## Restart all services
	@echo "$(BLUE)Restarting services...$(NC)"
	@docker-compose restart
	@echo "$(GREEN)✓$(NC) Services restarted"

restart-backend: ## Restart only backend
	@docker-compose restart backend
	@echo "$(GREEN)✓$(NC) Backend restarted"

restart-frontend: ## Restart only frontend
	@docker-compose restart frontend
	@echo "$(GREEN)✓$(NC) Frontend restarted"

shell-backend: ## Open shell in backend container
	@docker-compose exec backend sh

shell-frontend: ## Open shell in frontend container
	@docker-compose exec frontend sh

db-shell: ## Open PostgreSQL shell
	@docker-compose exec postgres psql -U postgres

db-backup: ## Backup PostgreSQL database
	@echo "$(BLUE)Creating database backup...$(NC)"
	@docker-compose exec postgres pg_dump -U postgres postgres > backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✓$(NC) Database backed up to backup_$$(date +%Y%m%d_%H%M%S).sql"

db-restore: ## Restore PostgreSQL database (make db-restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "$(RED)ERROR:$(NC) Please specify FILE=<backup_file.sql>"; \
		exit 1; \
	fi
	@echo "$(BLUE)Restoring database from $(FILE)...$(NC)"
	@docker-compose exec -T postgres psql -U postgres < $(FILE)
	@echo "$(GREEN)✓$(NC) Database restored"

migrate-up: ## Run database migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	@docker-compose exec backend /usr/local/bin/migrate -path=/app/migrations -database="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up
	@echo "$(GREEN)✓$(NC) Migrations complete"

migrate-down: ## Rollback last migration
	@echo "$(BLUE)Rolling back last migration...$(NC)"
	@docker-compose exec backend /usr/local/bin/migrate -path=/app/migrations -database="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" down 1
	@echo "$(GREEN)✓$(NC) Rollback complete"

prod: ## Start in production mode
	@echo "$(BLUE)Starting in production mode...$(NC)"
	@APP_ENV=production docker-compose up -d --build
	@echo "$(GREEN)✓$(NC) Production environment started"

stop: down ## Alias for 'down'

start: up ## Alias for 'up'
