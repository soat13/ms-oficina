.PHONY: install up down test mod fmt lint sh migrate-up migrate-down migrate-goto migrate-version

# Install and start the development environment
install:
	@if [ ! -f .env ]; then \
		cp .env-example .env; \
	fi
	docker compose down -v --remove-orphans
	docker compose up -d --build --force-recreate

# Up and down the development environment
up:
	@if [ ! -f .env ]; then \
		cp .env-example .env; \
	fi
	docker compose up -d
down:
	docker compose down

# Quality assurance commands
test:
	docker compose exec app-dev go test ./...
mod:
	docker compose exec app-dev go mod tidy
fmt:
	docker compose exec app-dev go fmt ./...
lint:
	docker compose exec app-dev go vet ./...

# Access the application container shell
sh:
	docker compose exec app-dev bash

# Database migrations
COMPOSE ?= docker compose --env-file .env

migrate-version:
	$(COMPOSE) run --rm --entrypoint sh migrate -lc 'migrate -source file:///migrations -database "$$DATABASE_URL" version'

migrate-up:
	$(COMPOSE) run --rm --entrypoint sh migrate -lc 'migrate -source file:///migrations -database "$$DATABASE_URL" up'

migrate-down:
	$(COMPOSE) run --rm --entrypoint sh migrate -lc 'migrate -source file:///migrations -database "$$DATABASE_URL" down 1'

migrate-goto:
	$(COMPOSE) run --rm --entrypoint sh migrate -lc 'migrate -source file:///migrations -database "$$DATABASE_URL" goto $(N)'

migrate-new:
	$(COMPOSE) run --rm --entrypoint sh migrate -lc 'migrate create -ext sql -dir /migrations -seq $(NAME)'
