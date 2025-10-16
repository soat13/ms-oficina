# Oficina API

## Sobre o Projeto

Sistema de gestão para oficinas mecânicas que automatiza o fluxo completo de atendimento, desde a entrada do veículo até a entrega ao cliente. A aplicação gerencia clientes, veículos, catálogos de serviços e produtos, orçamentos e ordens de serviço, proporcionando controle total sobre as operações da oficina.

### Funcionalidades Principais

- **Gestão de Clientes e Veículos** — Cadastro completo de clientes (donos dos veículos) e seus automóveis
- **Catálogos** — Gerenciamento de serviços técnicos (alinhamento, troca de óleo, etc.) e produtos (filtros, óleos, etc.) com controle de estoque
- **Orçamentos** — Criação de propostas detalhadas com itens de serviço e produtos, incluindo aprovação do cliente
- **Ordens de Serviço (OS)** — Fluxo completo desde a recepção do veículo até a liberação, passando por diagnóstico, aprovação, execução e finalização
- **Controle de Acesso** — Sistema de autenticação com diferentes níveis de permissão (Atendente, Mecânico, Gerente)
- **Eventos de Domínio** — Arquitetura orientada a eventos para desacoplar contextos (ex: redução de estoque ao aprovar orçamento)

A aplicação segue princípios de **Domain-Driven Design (DDD)** e **Arquitetura Hexagonal**, organizando o código em contextos de domínio bem definidos e isolados, facilitando manutenção e evolução.

---

## De-Para: Termos em Português → Inglês

### Contextos (Bounded Contexts)

| Português         | Inglês       | Pasta                    | Explicação                                                                     |
|-------------------|--------------|--------------------------|--------------------------------------------------------------------------------|
| Cliente           | Customer     | `internal/customer/`     | Representa o dono do veículo e titular do contrato                             |
| Usuário           | User         | `internal/user/`         | Usuários do sistema (Atendente, Mecânico, Gerente) com credenciais de acesso  |
| Veículo           | Vehicle      | `internal/vehicle/`      | Automóvel que será atendido na oficina                                         |
| Serviço           | Service      | `internal/service/`      | Atividade técnica do catálogo (ex: "Alinhamento", "Troca de Óleo")            |
| Produto           | Product      | `internal/product/`      | Peça ou insumo do catálogo (ex: "Filtro de Cabine", "Óleo 5W30")              |
| Ordem de Serviço  | Repair Order | `internal/repairorder/`  | OS - Autorização formal para executar serviços no veículo                      |
| Orçamento         | Estimate     | `internal/estimate/`     | Proposta de preços com itens de serviço e produtos                             |
| Autenticação      | Auth         | `internal/auth/`         | Sistema de login e geração de tokens JWT                                       |

### Atores (Roles)

| Português  | Inglês (código) | Constante   | Explicação                                                      |
|------------|-----------------|-------------|-----------------------------------------------------------------|
| Cliente    | Customer        | —           | Dono do veículo que solicita serviços e aprova orçamentos       |
| Atendente  | Attendant       | `attendant` | Cadastra clientes, cria OS, gerencia o fluxo de atendimento     |
| Mecânico   | Mechanic        | `mechanic`  | Executa diagnósticos e serviços técnicos                        |
| Gerente    | Manager         | `manager`   | Supervisiona operações e aprova ações administrativas           |

### Estados da Ordem de Serviço (Repair Order Status)

| Português            | Inglês (código)   | Constante            | Explicação                                 |
|----------------------|-------------------|----------------------|--------------------------------------------|
| Recebida             | Received          | `received`           | Veículo chegou na oficina                  |
| Em Diagnóstico       | In Diagnostics    | `in_diagnostics`     | Mecânico está avaliando o problema         |
| Aguardando Aprovação | Awaiting Approval | `awaiting_approval`  | Cliente precisa aprovar o orçamento        |
| Aprovada             | Approved          | `approved`           | Cliente aprovou, pode iniciar execução     |
| Em Execução          | In Execution      | `in_execution`       | Serviços sendo executados                  |
| Finalizada           | Finished          | `finished`           | Serviços concluídos, aguardando retirada   |
| Liberada             | Released          | `released`           | Veículo entregue ao cliente                |
| Cancelada            | Canceled          | `canceled`           | OS foi cancelada                           |

### Estados do Orçamento (Estimate Status)

| Português            | Inglês (código)   | Constante           | Explicação                                              |
|----------------------|-------------------|---------------------|---------------------------------------------------------|
| Aguardando Aprovação | Awaiting Approval | `awaiting_approval` | Orçamento criado, aguardando decisão do cliente         |
| Aguardando Estoque   | Awaiting Stock    | `awaiting_stock`    | Cliente aprovou, verificando disponibilidade de peças   |
| Aprovado             | Approved          | `approved`          | Totalmente aprovado e com estoque confirmado            |
| Rejeitado            | Rejected          | `rejected`          | Cliente recusou o orçamento                             |
| Cancelado            | Canceled          | `canceled`          | Orçamento foi cancelado (ex: OS cancelada)              |

---

## Aspectos Técnicos

API construída em **Go** utilizando o framework **Fiber** (HTTP) e **Bun** como ORM para PostgreSQL. A arquitetura é organizada por contextos de domínio independentes (`customer`, `user`, `vehicle`, `service`, `product`, `estimate`, `repairorder`, `auth`), seguindo os princípios de DDD e Arquitetura Hexagonal.

### Stack Tecnológica

- **Linguagem**: Go 1.25.1
- **Framework HTTP**: Fiber v2
- **Banco de Dados**: PostgreSQL
- **ORM**: Bun (query builder type-safe)
- **Migrações**: sql-migrate
- **Autenticação**: JWT (HS256)
- **Testes**: Go testing + testify
- **Containerização**: Docker + Docker Compose

### Arquitetura

Cada contexto de domínio é organizado em três camadas:

- **Domain** (`domain/`) — Entidades, value objects, regras de negócio e erros de domínio
- **Application** (`application/`) — Casos de uso (use cases), portas (interfaces) e lógica de aplicação
- **Infrastructure** (`infra/`) — Adaptadores concretos (repositórios Bun, handlers HTTP, serviços externos)

A comunicação entre contextos ocorre através de **eventos de domínio** gerenciados por um event bus em memória, mantendo o desacoplamento entre bounded contexts.

---

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

## Sem Make (comandos equivalentes)

- Subir os serviços:
  - `docker compose up -d`
- Aplicar migrações (instala a ferramenta dentro do container se necessário):
  - `docker compose exec -T app-dev sh -lc 'test -x /go/bin/sql-migrate || GOBIN=/go/bin /usr/local/go/bin/go install github.com/rubenv/sql-migrate/sql-migrate@latest'`
  - `docker compose exec -T app-dev sh -lc '/go/bin/sql-migrate up -config=./scripts/db/dbconfig.yml -env=development'`
- Executar a API:
  - `docker compose exec app-dev go run ./cmd/api`

## Variáveis de Ambiente

O arquivo `.env-example` contém todas as variáveis com valores de exemplo. Principais variáveis:

### Aplicação

- `PG_DSN` — String de conexão do PostgreSQL
- `PORT` — Porta onde a API escuta
- `JWT_SECRET` — Chave secreta para assinar tokens JWT (obrigatório)
- `JWT_EXPIRATION` — Tempo de expiração dos tokens (ex: `1h`, `90m`)

### PostgreSQL

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` — Credenciais do banco

### Testes

- `MIGRATIONS_DIR` — Diretório das migrações SQL
- `TEST_SEEDERS_DIR` — Diretório dos seeders de teste
- `DEBUG_TESTDB` — Ativa logs verbosos nos testes (valor: `1`)

### SonarQube

- `SONAR_JDBC_*` — Configurações de conexão com o banco do SonarQube
- `SONAR_TOKEN` — Token de autenticação para análise de código

### Outras

- `LOG_LEVEL` — Nível de log (`debug`, `info`, `warn`, `error`)
- `APP_ENV` — Ambiente de execução (`development`, `test`, `production`)

## Comandos úteis

### Gerenciamento de Containers

- `make up` — sobe os containers (DB e app-dev)
- `make down` — para os containers
- `make install` — reinstala todo ambiente (remove volumes, rebuilda containers)
- `make run` — executa a API no container

### Migrações de Banco

- `make migrate-up` — aplica todas as migrações pendentes
- `make migrate-down` — desfaz a última migração
- `make migrate-status` — exibe o status das migrações

## Estrutura principal

- `cmd/api/` — composição da aplicação (wiring)
- `internal/*/{domain,application,infra}` — camadas por contexto de domínio
- `scripts/db/migrations` — migrações SQL
- `tests/` — testes de integração e utilitários de banco

## Swagger / OpenAPI

Após subir a API, acesse a documentação:

- UI: `http://localhost/docs`
- Esquema: `http://localhost/openapi.yaml`

Observações:
- A UI usa assets do CDN (swagger-ui-dist). O arquivo do esquema fica embarcado no binário.
- As rotas documentadas correspondem às implementadas em todos os contextos.
 
## Autenticação JWT

A aplicação possui um fluxo de autenticação baseado em JSON Web Tokens para proteger as rotas administrativas (`/admin/**`). Os principais pontos são:

- **Login** — `POST /auth/login` recebe `email` e `password`, valida o usuário através do repositório em `internal/user/infra/db` (hash BCrypt) e retorna um JSON com `access_token`, `token_type` (`Bearer`), `expires_in`, `expires_at` e um resumo do usuário.
- **Geração/validação** — `internal/auth/infra/jwt` implementa `TokenService` usando HS256. As claims incluem `uid`, `email` e `roles`, com emissor `oficina-api` e expiração padrão de 1 hora.
- **Middleware** — `internal/auth/infra/http/middleware.go` valida o cabeçalho `Authorization: Bearer <token>` e anexa as claims ao contexto Fiber para uso posterior. Falhas retornam `401 Unauthorized` via o `ErrorHandler` compartilhado.
- **Bootstrap** — `internal/bootstrap/auth/setup.go` lê `JWT_SECRET` (obrigatório) e `JWT_EXPIRATION` (opcional, ex.: `90m`), registra o endpoint de login e instala o middleware antes das rotas `/admin`.

Toda a funcionalidade segue o padrão modular existente (`internal/auth/{domain,application,infra}`) e reutiliza a entidade `user` como fonte dos dados de acesso.

## Testes

- `make test` — executa todos os testes (unitários e de integração)
- `make test-coverage` — executa testes e gera relatório de coverage (`coverage.out`)
- `make sonar` — executa testes e envia análise para o SonarQube

**SonarQube:** Acesse http://localhost:9000 (login: `admin`/`admin`), gere um token em **My Account → Security → Generate Token** e adicione no `.env` como `SONAR_TOKEN`.

**Nota:** Os testes de integração criam bancos isolados automaticamente e aplicam migrações.
