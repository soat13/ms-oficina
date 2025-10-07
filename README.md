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

## Swagger / OpenAPI

Após subir a API, acesse a documentação:

- UI: `http://localhost/docs`
- Esquema: `http://localhost/openapi.yaml`

Observações:
- A UI usa assets do CDN (swagger-ui-dist). O arquivo do esquema fica embarcado no binário.
- As rotas documentadas correspondem às implementadas em `internal/service/infra/http` e `internal/estimate/infra/http`.
 
## Autenticação JWT

A aplicação possui um fluxo de autenticação baseado em JSON Web Tokens para proteger as rotas administrativas (`/admin/**`). Os principais pontos são:

- **Login** — `POST /auth/login` recebe `email` e `password`, valida o usuário através do repositório em `internal/user/infra/db` (hash BCrypt) e retorna um JSON com `access_token`, `token_type` (`Bearer`), `expires_in`, `expires_at` e um resumo do usuário.
- **Geração/validação** — `internal/auth/infra/jwt` implementa `TokenService` usando HS256. As claims incluem `uid`, `email` e `roles`, com emissor `oficina-api` e expiração padrão de 1 hora.
- **Middleware** — `internal/auth/infra/http/middleware.go` valida o cabeçalho `Authorization: Bearer <token>` e anexa as claims ao contexto Fiber para uso posterior. Falhas retornam `401 Unauthorized` via o `ErrorHandler` compartilhado.
- **Bootstrap** — `internal/bootstrap/auth/setup.go` lê `JWT_SECRET` (obrigatório) e `JWT_EXPIRATION` (opcional, ex.: `90m`), registra o endpoint de login e instala o middleware antes das rotas `/admin`.
- **Testes** — helpers em `tests/testsupport/auth.go` autenticam usuários reais durante os testes de integração, garantindo que as rotas protegidas só sejam acessadas com tokens válidos.

Toda a funcionalidade segue o padrão modular existente (`internal/auth/{domain,application,infra}`) e reutiliza a entidade `user` como fonte dos dados de acesso.
 

## Sem Make (comandos equivalentes)
- Subir os serviços:
  - `docker compose up -d`
- Aplicar migrações (instala a ferramenta dentro do container se necessário):
  - `docker compose exec -T app-dev sh -lc 'test -x /go/bin/sql-migrate || GOBIN=/go/bin /usr/local/go/bin/go install github.com/rubenv/sql-migrate/sql-migrate@latest'`
  - `docker compose exec -T app-dev sh -lc '/go/bin/sql-migrate up -config=./scripts/db/dbconfig.yml -env=development'`
- Executar a API:
  - `docker compose exec app-dev go run ./cmd/api`

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

### Executar Testes
- Rodar todos os testes:
  ```bash
  make test
  ```

- Rodar testes com coverage completo:
  ```bash
  make test-coverage
  ```
  Este comando:
  - Executa todos os testes (unitários e de integração)
  - Calcula coverage de **todos** os pacotes em `internal/` e `pkg/`
  - Gera o arquivo `coverage.out`
  - Os testes de integração contam para o coverage dos use cases

### SonarQube
Para executar análise de código e enviar coverage para o SonarQube:

```bash
make sonar
```

Este comando:
- Executa todos os testes com coverage completo
- Gera o relatório de coverage (`coverage.out`)
- Envia a análise para o SonarQube (http://localhost:9000)

**Observações:**
- Os testes de integração sobem um banco isolado por teste e aplicam migrações automaticamente
- O coverage inclui código testado pelos testes de integração em `tests/integration/`
- O SonarQube deve estar rodando localmente na porta 9000

## Comandos úteis
- `make migrate-status` — status das migrações
- `make migrate-down` — desfaz a última migração
- `make fmt` — formata o código
- `make lint` — `go vet`
- `make mock` — gera mocks para testes
- `make sh` — shell no container `app-dev`

## Estrutura principal
- `cmd/api/` — composição da aplicação (wiring)
- `internal/*/{domain,application,infra}` — camadas por contexto de domínio
- `scripts/db/migrations` — migrações SQL
- `tests/` — testes de integração e utilitários de banco
