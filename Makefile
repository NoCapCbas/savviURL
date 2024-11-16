
make prod:
	docker compose --env-file .env.prod -f docker-compose.prod.yml up --build -d

make prod-down:
	docker compose --env-file .env.prod -f docker-compose.prod.yml down

make prod-clean:
	make prod-down
	docker system prune -f
	docker volume rm url-shortening-service_url-shortener-redis-data

make dev:
	docker compose --env-file .env.dev -f docker-compose.dev.yml up --build -d

make dev-down:
	docker compose --env-file .env.dev -f docker-compose.dev.yml down

make dev-clean:
	make dev-down
	docker system prune -f
	docker volume rm url-shortening-service_redis_data

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  prod            Run the production environment"
	@echo "  prod-down       Stop the production environment"
	@echo "  prod-clean      Clean up the production environment"
	@echo "  dev             Run the development environment"
	@echo "  dev-down        Stop the development environment"
	@echo "  dev-clean       Clean up the development environment"
	@echo "  help            Display this help message"

.PHONY: prod prod-down prod-clean dev dev-down dev-clean help