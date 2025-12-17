# Oficina API

## Sobre o Projeto

Sistema de gestão para oficinas mecânicas que automatiza o fluxo completo de atendimento, desde a entrada do veículo até a entrega ao cliente. A aplicação gerencia clientes, veículos, catálogos de serviços e produtos, orçamentos e ordens de serviço, oferecendo controle centralizado das operações da oficina.

---

## Funcionalidades Principais

- **Gestão de Clientes e Veículos** — Cadastro e relacionamento entre clientes e seus automóveis
- **Catálogos** — Serviços técnicos e produtos com controle de estoque
- **Orçamentos** — Propostas com itens de serviço e produto, sujeitas à aprovação
- **Ordens de Serviço (OS)** — Fluxo completo: recepção, diagnóstico, aprovação, execução e liberação
- **Controle de Acesso** — Autenticação e autorização por perfis (Atendente, Mecânico, Gerente)
- **Eventos de Domínio** — Desacoplamento entre contextos (ex: baixa de estoque após aprovação)

A aplicação segue os princípios de **Domain-Driven Design (DDD)** e **Arquitetura Hexagonal**, com organização por contextos de domínio independentes.

---

## Fluxos Principais

- **Repair Order — fluxo end-to-end**  
  Fluxo operacional completo da Ordem de Serviço, da criação à entrega do veículo.  
   [`docs/repair-order-complete-flow.md`](docs/repair-order-complete-flow.md)

> Payloads, schemas e exemplos de request/response estão documentados no **Swagger/OpenAPI**.  
> Os documentos de fluxo focam exclusivamente na **ordem de chamadas e regras de negócio**.

---

## Documentação Técnica

- Arquitetura e decisões técnicas: [`docs/architecture.md`](docs/architecture.md)
- Glossário de domínio (termos, roles e status): [`docs/domain-glossary.md`](docs/domain-glossary.md)

---

## Aspectos Técnicos

A API é construída em **Go**, seguindo princípios de **DDD** e **Arquitetura Hexagonal**, organizada por contextos de domínio independentes.

Diagramas e detalhes arquiteturais: [`docs/architecture.md`](docs/architecture.md)

### Stack Tecnológica

- **Linguagem**: Go 1.25.1
- **Framework HTTP**: Fiber v2
- **ORM**: Bun
- **Migrações**: sql-migrate
- **Autenticação**: JWT (HS256)
- **Testes**: Go testing + testify
- **Containerização**: Docker + Docker Compose
- **Banco de Dados**: PostgreSQL — Escolhido por ser um SGBD relacional maduro e confiável, adequado para garantir integridade transacional em operações críticas como criação de Ordens de Serviço, aprovação de orçamentos e controle de estoque. O modelo relacional facilita a consistência entre entidades fortemente relacionadas e oferece suporte nativo a transações ACID, constraints e índices, essenciais para a confiabilidade e evolução do sistema.

---

## Requisitos

- Docker e Docker Compose
- make (opcional, recomendado)

---

## Executando o Projeto

### Com Make

1. Subir os containers:
    - `make install`
    - Caso não exista `.env`, ele será criado a partir do `.env-example`
2. Aplicar migrações:
    - `make migrate-up`
3. Rodar a API:
    - `make run`

### Sem Make

1. Criar o `.env` a partir do `.env-example`
2. Subir os containers:
    - `docker compose up -d`
3. Aplicar migrações:
    - `docker compose exec -T app-dev sh -lc 'test -x /go/bin/sql-migrate || GOBIN=/go/bin /usr/local/go/bin/go install github.com/rubenv/sql-migrate/sql-migrate@latest'`
    - `docker compose exec -T app-dev sh -lc '/go/bin/sql-migrate up -config=./scripts/db/dbconfig.yml -env=development'`
4. Rodar a API:
    - `docker compose exec app-dev go run ./cmd/api`

---

## Swagger / OpenAPI

- UI: http://localhost:8080/docs
- Spec: http://localhost:8080/openapi.yaml
- Arquivo no repositório: `assets/docs/openapi.yaml`

---

## Autenticação

A API utiliza **JWT** para proteger as rotas administrativas (`/admin/**`).

- Login: `POST /auth/login`
- Header: `Authorization: Bearer <token>`

---

## Testes

- `make test` — executa testes unitários e de integração
- `make test-coverage` — gera relatório de cobertura
- `make sonar` — executa análise no SonarQube

> Os testes de integração criam bancos isolados e aplicam migrações automaticamente.

---

## Links Úteis

- GitHub: https://github.com/soat13/fase-1-oficina
- Miro: https://miro.com/app/board/uXjVJLyIcr8=
