.PHONY: install up down run test test-bdd test-coverage mod vendor tidy fmt lint sh \
        migrate-up migrate-down migrate-status seed-up sonar sonar-analysis

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Infra
up:
	@if [ ! -f .env ]; then cp .env-example .env; fi
	docker compose up -d
	@docker compose up -d --wait app-dev

down:
	docker compose down

install:
	@if [ ! -f .env ]; then cp .env-example .env; fi
	docker compose down -v --remove-orphans
	docker compose pull
	docker compose up -d --build --force-recreate
	@docker compose up -d --wait app-dev

run:
	docker compose exec app-dev go run ./cmd/api

# QA
test:
	docker compose exec app-dev go test ./...
test-bdd:
	docker compose exec app-dev go test ./tests/bdd/... -v
test-coverage:
	docker compose exec app-dev go test -coverpkg=./internal/... -coverprofile=coverage.out ./...
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
    	  -source=internal/auth/application/ports.go \
    	  -destination=internal/auth/application/mocks/ports_mock.go \
    	  -package=mocks'

# -------------------------------
# SonarQube Analysis
# -------------------------------

sonar: up
	@if [ -z "$(SONAR_TOKEN)" ]; then \
		echo "Error: SONAR_TOKEN is not set. Please set it in your .env file."; \
		exit 1; \
	fi
	make test-coverage
	docker compose run --rm sonar-scanner \
		-Dsonar.projectKey=oficina \
		-Dsonar.projectName="Fase 1 - Oficina" \
		-Dsonar.sources=. \
		-Dsonar.exclusions=**/vendor/**,**/mocks/**,**/*_test.go,**/tests/**,**/scripts/**,**/assets/**,.env*,**/*.md,**/cmd/**,**/views.go,**/repository.go,**/bootstrap/** \
		-Dsonar.tests=. \
		-Dsonar.test.inclusions=**/*_test.go \
		-Dsonar.go.coverage.reportPaths=coverage.out