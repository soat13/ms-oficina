up:
	docker compose up -d
down:
	docker compose down -v

run:
	docker compose exec dev go run ./cmd/api
test:
	docker compose exec dev go test ./...

mod:
	docker compose exec dev go mod tidy
fmt:
	docker compose exec dev go fmt ./...
lint:
	docker compose exec dev go vet ./...

sh:
	docker compose exec dev bash
