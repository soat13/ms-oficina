.PHONY: install up down run test mod vendor tidy fmt lint sh \
        migrate-up migrate-down migrate-status seed-up

# Infra
up:
	@if [ ! -f .env ]; then cp .env-example .env; fi
	docker compose up -d

down:
	docker compose down

install:
	@if [ ! -f .env ]; then cp .env-example .env; fi
	docker compose down -v --remove-orphans
	docker compose pull
	docker compose up -d --build --force-recreate

run:
	docker compose exec app-dev go run ./cmd/api/main.go

# QA
test:
	docker compose exec app-dev go test ./...
mod:
	docker compose exec app-dev go mod tidy
vendor:
	docker compose exec app-dev go mod vendor
tidy:
	docker compose exec app-dev go mod tidy
fmt:
	docker compose exec app-dev go fmt ./...
lint:
	docker compose exec app-dev go vet ./...

# Shell
sh:
	docker compose exec app-dev bash

# -------------------------------
# sql-migrate running on app-dev
# -------------------------------
COMPOSE ?= docker compose --env-file .env
GO_BIN           ?= /usr/local/go/bin/go
SQL_MIGRATE_BIN  ?= /go/bin/sql-migrate
SQL_MIGRATE_PKG  ?= github.com/rubenv/sql-migrate/sql-migrate@latest
SQL_MIGRATE_CFG  ?= ./scripts/db/dbconfig.yml

migrate-install: up
	$(COMPOSE) exec -T app-dev sh -lc 'test -x $(SQL_MIGRATE_BIN) || GOBIN=/go/bin $(GO_BIN) install $(SQL_MIGRATE_PKG)'

migrate-status: up migrate-install
	$(COMPOSE) exec -T app-dev sh -lc '$(SQL_MIGRATE_BIN) status -config=$(SQL_MIGRATE_CFG) -env=development'

migrate-up: up migrate-install
	$(COMPOSE) exec -T app-dev sh -lc '$(SQL_MIGRATE_BIN) up -config=$(SQL_MIGRATE_CFG) -env=development'

migrate-down: up migrate-install
	$(COMPOSE) exec -T app-dev sh -lc '$(SQL_MIGRATE_BIN) down -config=$(SQL_MIGRATE_CFG) -env=development -limit=1'

# -------------------------------
# Mockgen running on app-dev
# -------------------------------
DOCKER_EXEC := docker compose exec app-dev
GO_BIN      := /usr/local/go/bin/go
MOCKGEN     := /go/bin/mockgen

mockgen-install:
	$(DOCKER_EXEC) sh -lc 'test -x $(MOCKGEN) || $(GO_BIN) install github.com/golang/mock/mockgen@latest'

mock: mockgen-install
	$(DOCKER_EXEC) sh -lc '$(MOCKGEN) \
	  -source=internal/repairorder/application/repository.go \
	  -destination=internal/repairorder/mocks/repository_mock.go \
	  -package=mocks'
