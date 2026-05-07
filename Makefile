# Transcendence Makefile

.PHONY: help build up down logs stop restart clean test dev prod migrate seed dev-backend logs-backend-only shell

NAME = Transcendence
DOCKER_COMPOSE = docker compose -f infra/docker-compose.yml
DOCKER_COMPOSE_PROD = docker compose -f infra/docker-compose.prod.yml

SERVICES_BACKEND = backend db

# ==================== Main Commands ====================

help:
	@echo "╔═════════════════════════════════════════════╗"
	@echo "║    $(NAME) - Development Commands     ║"
	@echo "╚═════════════════════════════════════════════╝"
	@echo ""
	@echo "Quick Start:"
	@echo "  make dev                - Full dev setup (build + up + logs)"
	@echo "  make dev-backend        - Start backend + db only"
	@echo "  make re                 - Full restart (down + clean + build + up)"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build              - Build all Docker images"
	@echo "  make up                 - Start all containers"
	@echo "  make down               - Stop containers (keep volumes)"
	@echo "  make stop               - Stop containers gracefully"
	@echo "  make restart            - Restart containers"
	@echo "  make clean              - Remove containers and volumes"
	@echo "  make ps                 - Show container status"
	@echo ""
	@echo "Logs & Debugging:"
	@echo "  make logs               - Show all logs (follow mode)"
	@echo "  make logs-backend       - Show backend logs only"
	@echo "  make logs-backend-only  - Backend + DB logs"
	@echo "  make logs-frontend      - Show frontend logs only"
	@echo "  make logs-nginx         - Show nginx logs only"
	@echo "  make logs-db            - Show database logs only"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run all tests"
	@echo "  make test-backend       - Run backend tests"
	@echo "  make test-frontend      - Run frontend tests"
	@echo ""
	@echo "Database:"
	@echo "  make migrate            - Run database migrations"
	@echo "  make seed               - Seed database with test data"
	@echo ""
	@echo "Development:"
	@echo "  make shell              - Access backend container shell"
	@echo "  make shell-backend      - Access backend container shell"
	@echo "  make shell-frontend     - Access frontend container shell"
	@echo "  make shell-db           - Access database shell"
	@echo "  make health             - Check health of all services"
	@echo "  make fmt                - Format Go code"
	@echo "  make lint               - Lint Go code"
	@echo ""
	@echo "Production:"
	@echo "  make build-prod         - Build production images"
	@echo "  make up-prod            - Start production containers"
	@echo "  make down-prod          - Stop production containers"
	@echo ""
	@echo "Utilities:"
	@echo "  make prune              - Clean up Docker resources"
	@echo "  make version            - Show versions"
	@echo ""

# ==================== Quick Commands ====================

dev: build up logs
	@echo ""
	@echo "Dev environment ready!"

dev-backend:
	@echo "Starting backend + database only..."
	@$(DOCKER_COMPOSE) up -d --build $(SERVICES_BACKEND)
	@echo ""
	@echo "╔════════════════════════════════════════╗"
	@echo "║     Backend Dev Started!               ║"
	@echo "╚════════════════════════════════════════╝"
	@echo ""
	@echo "Backend:   http://localhost:8080"
	@echo "Database:  postgres://localhost:5432"

all: build up
	@echo ""
	@echo "All services started!"

# ==================== Docker Commands ====================

build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build
	@echo "Build complete!"

up:
	@echo "Starting containers..."
	@$(DOCKER_COMPOSE) up -d
	@sleep 2
	@echo ""
	@echo "╔════════════════════════════════════════╗"
	@echo "║     Services Started!                  ║"
	@echo "╚════════════════════════════════════════╝"
	@echo ""
	@echo "Frontend:  http://localhost"
	@echo "Backend:   http://localhost/api"
	@echo "Database:  postgres://localhost:5432"
	@echo ""

down:
	@echo "Stopping containers..."
	@$(DOCKER_COMPOSE) down
	@echo "Services stopped"

stop:
	@echo "Stopping containers gracefully..."
	@$(DOCKER_COMPOSE) stop
	@echo "Services stopped"

restart: down up
	@echo "Restart complete!"

re:
	@echo "Stopping containers..."
	@make down
	@echo "Removing backend images..."
	@docker rmi $$(docker images | grep backend | awk '{print $$3}') || true
	@echo "Starting again..."
	@make up

clean:
	@echo "Cleaning up containers and volumes..."
	@$(DOCKER_COMPOSE) down -v
	@echo "Cleanup complete!"

ps:
	@echo "Container status:"
	@$(DOCKER_COMPOSE) ps

# ==================== Logs ====================

logs:
	@echo "Showing all logs (Press Ctrl+C to stop)..."
	@$(DOCKER_COMPOSE) logs -f

logs-backend:
	@echo "Backend logs:"
	@$(DOCKER_COMPOSE) logs -f backend

logs-backend-only:
	@echo "Backend + DB logs:"
	@$(DOCKER_COMPOSE) logs -f backend db

logs-frontend:
	@echo "Frontend logs:"
	@$(DOCKER_COMPOSE) logs -f frontend

logs-nginx:
	@echo "Nginx logs:"
	@$(DOCKER_COMPOSE) logs -f nginx

logs-db:
	@echo "Database logs:"
	@$(DOCKER_COMPOSE) logs -f db

# ==================== Testing ====================

test:
	@echo "Running all tests..."
	@$(DOCKER_COMPOSE) exec -T backend go test ./...
	@echo "Tests complete!"

test-backend:
	@echo "Running backend tests locally..."
	@cd backend && DB_HOST=localhost go test ./tests/... -v -count=1 2>&1 | \
		sed 's/--- PASS/\x1b[32m--- PASS\x1b[0m/g' | \
		sed 's/--- FAIL/\x1b[31m--- FAIL\x1b[0m/g' | \
		sed 's/^ok/\x1b[32mok\x1b[0m/g' | \
		sed 's/^FAIL/\x1b[31mFAIL\x1b[0m/g'

# test-backend:
# 	@echo "Running backend tests..."
# 	@$(DOCKER_COMPOSE) exec -T backend go test ./...

test-frontend:
	@echo "Running frontend tests..."
	@$(DOCKER_COMPOSE) exec -T frontend npm test

# ==================== Database ====================

seed:
	@echo "Seeding database via Docker..."
	@$(DOCKER_COMPOSE) --profile seed run --rm seed
	@echo "Seed complete!"

seed-clean: clean up
	@echo "Fresh DB, now seeding..."
	@sleep 3
	@$(DOCKER_COMPOSE) --profile seed run --rm seed
	@echo "Clean seed complete!"

# ==================== Shell Access ====================

shell:
	@echo "Accessing backend container..."
	@$(DOCKER_COMPOSE) exec backend sh

shell-backend:
	@echo "Accessing backend container..."
	@$(DOCKER_COMPOSE) exec backend sh

shell-frontend:
	@echo "Accessing frontend container..."
	@$(DOCKER_COMPOSE) exec frontend sh

shell-db:
	@echo "Accessing database container..."
	@$(DOCKER_COMPOSE) exec db psql -U app -d app_db

# ==================== Health Checks ====================

health:
	@echo "Checking service health..."
	@echo ""
	@echo "Backend:  $$(curl -s http://localhost:8080/health || echo 'Down')"
	@echo "Frontend: $$(curl -s http://localhost:3000 | head -1 || echo 'Down')"
	@echo "Nginx:    $$(curl -s http://localhost | head -1 || echo 'Down')"
	@echo ""
	@$(DOCKER_COMPOSE) ps

# ==================== Development Tools ====================

fmt:
	@echo "Formatting code..."
	@cd backend && go fmt ./...
	@echo "Format complete!"

lint:
	@echo "Linting code..."
	@cd backend && go vet ./...
	@echo "Lint complete!"

# ==================== Production Commands ====================

build-prod:
	@echo "Building production images..."
	@$(DOCKER_COMPOSE_PROD) build
	@echo "Production build complete!"

up-prod:
	@echo "Starting production containers..."
	@$(DOCKER_COMPOSE_PROD) up -d
	@echo "Production services started!"

down-prod:
	@echo "Stopping production containers..."
	@$(DOCKER_COMPOSE_PROD) down
	@echo "Production services stopped"

# ==================== Utils ====================

prune:
	@echo "Pruning unused Docker resources..."
	@docker system prune -f
	@echo "Prune complete!"

version:
	@echo "$(NAME) Makefile v1.0"
	@echo "Docker: $$(docker --version)"
	@echo "Compose: $$(docker compose --version)"

.DEFAULT_GOAL := help
