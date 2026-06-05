# ft_transcendence Makefile
#
# One environment, one compose file. What changes between machines/deploys is
# only infra/.env (compose reads it; fixed topology like DB_HOST/REDIS_HOST lives
# in infra/docker-compose.yml). There is no dev/prod split here.
#
# Container engine: Podman is preferred (the school runs Fedora); Docker is used
# as a fallback. Both are supported — override the auto-detection with:
#   make ENGINE=docker up
# If `podman compose` is unavailable on your box, point COMPOSE at a provider:
#   make COMPOSE="podman-compose -f infra/docker-compose.yml" up

NAME       := ft_transcendence
ENGINE     ?= $(shell command -v podman >/dev/null 2>&1 && echo podman || echo docker)
COMPOSE    ?= $(ENGINE) compose -f infra/docker-compose.yml
GO_IMAGE   := golang:1.25

.PHONY: help build up down stop restart re clean ps certs \
        logs logs-backend logs-frontend logs-nginx logs-db \
        test test-frontend seed \
        shell-backend shell-frontend shell-nginx shell-postgres shell-redis \
        fmt lint swagger prune version

CERT_DIR := infra/certs
.DEFAULT_GOAL := help

# ==================== Help ====================

help:
	@echo "$(NAME) — using container engine: $(ENGINE)"
	@echo ""
	@echo "Stack:"
	@echo "  make up            Build and start everything (detached)"
	@echo "  make down          Stop containers (keep volumes)"
	@echo "  make stop          Stop containers without removing them"
	@echo "  make restart       down + up"
	@echo "  make re            Rebuild from clean state (down -v + build + up)"
	@echo "  make clean         Stop and remove containers + volumes"
	@echo "  make ps            Show container status"
	@echo "  make certs         Generate the dev TLS cert once (auto-run by 'up')"
	@echo ""
	@echo "Logs (follow):"
	@echo "  make logs          All services      make logs-backend   Backend"
	@echo "  make logs-frontend Frontend          make logs-nginx     Proxy (nginx)"
	@echo "  make logs-db       Database"
	@echo ""
	@echo "Tests:"
	@echo "  make test          Backend tests inside a $(GO_IMAGE) container"
	@echo "  make test-frontend Frontend tests (inside the frontend container)"
	@echo ""
	@echo "Database:"
	@echo "  make seed          Seed via the REST API with Python (500 users,"
	@echo "                     2000 unique posts, comments, likes, follows;"
	@echo "                     real photos). Override: make seed USERS=200 POSTS_TARGET=800"
	@echo ""
	@echo "Shells:"
	@echo "  make shell-backend   Open a shell in the backend (Go API) container"
	@echo "  make shell-frontend  Open a shell in the frontend (web) container"
	@echo "  make shell-nginx     Open a shell in the nginx reverse-proxy container"
	@echo "  make shell-postgres  Open a psql prompt in the PostgreSQL container"
	@echo "  make shell-redis     Open a redis-cli prompt in the Redis container"
	@echo ""
	@echo "Tools:"
	@echo "  make fmt             Format Go code (go fmt ./...)"
	@echo "  make lint            Vet Go code (go vet ./...)"
	@echo "  make swagger         Regenerate OpenAPI spec (backend/docs)"
	@echo "  make prune           Reclaim unused engine resources"
	@echo "  make version         Show engine and compose versions"
	@echo ""
	@echo "Override the engine: make ENGINE=docker <target>"

# ==================== Stack ====================

up: certs
	@echo "Starting stack with $(ENGINE)..."
	@$(COMPOSE) up -d --build
	@echo "Up. App: https://localhost:3000 (self-signed cert). 'make ps' for ports, 'make logs' to follow."
	@echo "API docs (Swagger): https://localhost:3000/swagger/index.html"

down:
	@$(COMPOSE) down

stop:
	@$(COMPOSE) stop

restart: down up

re: clean up

clean:
	@$(COMPOSE) down -v

ps:
	@$(COMPOSE) ps

# Generate a self-signed cert ONCE into infra/certs/ (gitignored), bind-mounted
# into the proxy. Stable across rebuilds, so the browser "accept the risk"
# exception you grant once keeps working. Delete the files to force a new cert.
# `up` depends on this, so a fresh checkout gets a cert automatically.
certs:
	@if [ -f $(CERT_DIR)/dev.crt ] && [ -f $(CERT_DIR)/dev.key ]; then \
		echo "TLS cert already present in $(CERT_DIR)/ — skipping (rm them to regenerate)."; \
	else \
		echo "Generating self-signed dev TLS cert into $(CERT_DIR)/ ..."; \
		mkdir -p $(CERT_DIR); \
		openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 \
			-keyout $(CERT_DIR)/dev.key -out $(CERT_DIR)/dev.crt \
			-subj "/C=FR/ST=Ile-de-France/L=Paris/O=42/OU=ft_transcendence/CN=ft_transcendence.42.fr" \
			-addext "subjectAltName=DNS:ft_transcendence.42.fr,DNS:localhost,IP:127.0.0.1" >/dev/null 2>&1; \
		echo "Done. On first visit to https://localhost:3000, accept the warning once — it now persists."; \
	fi

# ==================== Logs ====================

logs:
	@$(COMPOSE) logs -f

logs-backend:
	@$(COMPOSE) logs -f backend

logs-frontend:
	@$(COMPOSE) logs -f frontend

logs-nginx:
	@$(COMPOSE) logs -f proxy

logs-db:
	@$(COMPOSE) logs -f db

# ==================== Tests ====================

# Backend tests run inside a Go container. The suite uses testcontainers, so the
# container talks to the host engine through its socket to spin up Postgres.
# Podman (rootless) needs its API socket enabled once:
#   systemctl --user enable --now podman.socket
test:
	@echo "Running backend tests in $(GO_IMAGE) via $(ENGINE)..."
	@sock=$$( [ "$(ENGINE)" = "podman" ] \
		&& echo "$${XDG_RUNTIME_DIR:-/run/user/$$(id -u)}/podman/podman.sock" \
		|| echo /var/run/docker.sock ); \
	$(ENGINE) run --rm \
		-v "$(CURDIR)/backend":/app:z -w /app \
		-v "$$sock":/var/run/docker.sock:z \
		-v transcendence-gomod:/go/pkg/mod \
		-e DOCKER_HOST=unix:///var/run/docker.sock \
		-e TESTCONTAINERS_RYUK_DISABLED=true \
		--env-file infra/.env \
		$(GO_IMAGE) go test ./test/... -count=1

test-frontend:
	@$(COMPOSE) exec -T frontend npm test

# ==================== Database ====================

# Seed through the public REST API with Python (stdlib only). Pulls real names
# and photos from randomuser.me, then creates users/posts/comments/likes/follows.
# Runs on the host, so it needs python3 and a running stack (the script waits for
# the API to come up). The seeder paces its auth calls to RATE_LIMIT_MAX (the
# backend's per-IP/min ceiling) — by default we read that straight from infra/.env
# so the seeder matches the running backend. Override the user count with USERS=N,
# or the limit with RATE_LIMIT_MAX=N on the command line.
RATE_LIMIT_MAX ?= $(shell sed -n 's/^RATE_LIMIT_MAX=//p' infra/.env 2>/dev/null | tail -1)
seed:
	@$(if $(USERS),USERS="$(USERS)") $(if $(POSTS_TARGET),POSTS_TARGET="$(POSTS_TARGET)") $(if $(PAR),PAR="$(PAR)") \
		$(if $(RATE_LIMIT_MAX),RATE_LIMIT_MAX="$(RATE_LIMIT_MAX)") \
		python3 scripts/seed.py

# ==================== Shells & tools ====================

shell-backend:
	@$(COMPOSE) exec backend sh

shell-frontend:
	@$(COMPOSE) exec frontend sh

shell-nginx:
	@$(COMPOSE) exec proxy sh

shell-postgres:
	@$(COMPOSE) exec db psql -U app -d app_db

shell-redis:
	@$(COMPOSE) exec redis redis-cli

fmt:
	@cd backend && go fmt ./...

lint: swagger
	@cd backend && go vet ./...

swagger:
	@echo "Regenerating OpenAPI spec into backend/docs ..."
	@cd backend && go tool swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Done. Start the stack ('make up') and open https://localhost:3000/swagger/index.html"

# ==================== Utils ====================

prune:
	@echo "Reclaiming unused $(ENGINE) data (containers, networks, images, cache; named volumes kept)..."
	$(ENGINE) system prune -f

version:
	@echo "$(NAME) — engine: $(ENGINE)"
	@$(ENGINE) --version
	@$(COMPOSE) version
