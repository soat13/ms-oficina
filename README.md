# Oficina API — Guia Rápido

API em Go (Fiber + Bun/Postgres) organizada por contextos (DDD/Hexagonal):
`service`, `estimate` e `repairorder`. Este guia mostra como subir o ambiente,
rodar migrações, executar a API e rodar testes.

## Requisitos
- Docker e Docker Compose
- make (opcional, recomendado)

## Subir o ambiente (com Make)
1. Suba os containers (DB e app-dev):
   - `make up`
   - Dica: se não existir `.env`, ele será criado a partir de `.env-example`.
2. Aplique as migrações de banco:
   - `make migrate-up`
3. Rode a API:
   - `make run`

Por padrão a aplicação escuta na porta `8080` dentro do container e está
exposta no host em `http://localhost` (porta 80 mapeada para 8080).

## TODO

 - CRUD: product (Pisani)
 - CRUD: customer (Lucas)
 - CRUD: vehicles (Marcos)
 - Criar swagger (Pisani)
 - Autenticação JWT
  - Verificar nome da tabela com a linguagem oblíqua

## Sem Make (comandos equivalentes)
- Subir os serviços:
  - `docker compose up -d`
- Aplicar migrações (instala a ferramenta dentro do container se necessário):
  - `docker compose exec -T app-dev sh -lc 'test -x /go/bin/sql-migrate || GOBIN=/go/bin /usr/local/go/bin/go install github.com/rubenv/sql-migrate/sql-migrate@latest'`
  - `docker compose exec -T app-dev sh -lc '/go/bin/sql-migrate up -config=./scripts/db/dbconfig.yml -env=development'`
- Executar a API:
  - `docker compose exec app-dev go run ./cmd/api/main.go`

## Variáveis de ambiente
- `PG_DSN`: string de conexão do Postgres (definida no `.env` e no compose)
- `PORT`: porta interna da API (default 8080 dentro do container)

## Endpoints úteis (Admin Services)
Base URL (host): `http://localhost`

- Criar serviço:
  - `POST /admin/services/`
  - Body JSON:
    `{ "name": "Alignment", "price_cents": 12000, "currency": "BRL" }`

- Listar serviços:
  - `GET /admin/services/?limit=50&offset=0`

- Buscar por ID:
  - `GET /admin/services/{id}`

- Atualizar serviço:
  - `PUT /admin/services/{id}`
  - Body JSON (parcial): `{ "name": "New Name" }`

- Deletar serviço:
  - `DELETE /admin/services/{id}`

Observação: o endpoint de orçamento (`POST /repair-orders/{id}/estimate`) requer
uma ordem de reparo e catálogos válidos no banco. O repositório ainda não expõe
rotas HTTP para criar ordens de reparo; utilize inserções no banco ou os testes
de integração como referência.

## Testes
- Rodar todos os testes:
  - `make test`

Os testes de integração sobem um banco isolado por teste, aplicam migrações e
seeds de teste automaticamente.

## Comandos úteis
- `make migrate-status` — status das migrações
- `make migrate-down` — desfaz a última migração
- `make fmt` — formata o código
- `make lint` — `go vet`
- `make sh` — shell no container `app-dev`

## Estrutura principal
- `cmd/api/` — composição da aplicação (wiring)
- `internal/*/{domain,application,infra}` — camadas por contexto de domínio
- `scripts/db/migrations` — migrações SQL
- `tests/` — testes de integração e utilitários de banco
