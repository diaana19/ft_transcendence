NAME = Transcendence

all: build up

re: down clean build up

help:
	@echo "=== Transcendence Project ==="
	@echo ""
	@echo "Available commands:"
	@echo "  make build              - Build Docker images"
	@echo "  make up                 - Start containers"
	@echo "  make down               - Stop containers"
	@echo "  make logs               - Display logs"
	@echo "  make stop               - Stop containers"
	@echo "  make restart            - Restart containers"
	@echo "  make clean              - Remove containers and volumes"
	@echo ""

build:
	docker compose build

up:
	docker compose up -d
	@echo "✅ Services started!"
	@echo ✅ backend started!
	@echo ✅ frontend started!

down:
	docker compose down -v
	@echo "✅ Services stopped"

logs:
	docker compose logs -f

stop:
	docker compose stop
	@echo "✅ Services stopped"

restart: down up

clean:
	docker compose down -v
	@echo "✅ Containers and volumes removed"
