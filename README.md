# Oficina API

## Sumário

- [Sobre o Projeto](#sobre-o-projeto)
- [Fluxos Principais](#fluxos-principais)
- [Documentação Técnica](#documentação-técnica)
- [Aspectos Técnicos](#aspectos-técnicos)
    - [Arquitetura](#arquitetura)
    - [Stack Tecnológica](#stack-tecnológica)
- [Requisitos](#requisitos)
- [Pipeline](#pipeline)
- [Executando o Projeto](#executando-o-projeto)
- [Swagger / OpenAPI](#swagger--openapi)
- [Testes](#testes)
- [Deploy no Kubernetes](#deploy-no-kubernetes)
- [Arquitetura de Eventos (SQS + SNS)](#arquitetura-de-eventos-sqs--sns)
- [SAGA Pattern e Compensações](#saga-pattern-e-compensações)
- [Observabilidade](#observabilidade)

## Sobre o Projeto

Sistema de gestão para oficinas mecânicas que automatiza o fluxo completo de atendimento, desde a entrada do veículo até a entrega ao cliente. A aplicação gerencia clientes, veículos, catálogos de serviços e produtos, orçamentos e ordens de serviço, oferecendo controle centralizado das operações da oficina.

### Links Úteis

- GitHub: https://github.com/soat13/oficina
- Miro: https://miro.com/app/board/uXjVJLyIcr8=

### Funcionalidades Principais

- **Gestão de Clientes e Veículos** — Cadastro e relacionamento entre clientes e seus automóveis
- **Catálogos** — Serviços técnicos e produtos com controle de estoque
- **Orçamentos** — Propostas com itens de serviço e produto, sujeitas à aprovação
- **Ordens de Serviço (OS)** — Fluxo completo: recepção, diagnóstico, aprovação, execução e liberação
- **Pagamentos** — Integração via eventos com o serviço `payment` para atualização do status da OS conforme mudanças de status de pagamento
- **Eventos de Domínio via SQS + SNS** — Comunicação assíncrona entre contextos: filas SQS para entrega ponto a ponto e tópicos SNS para fan-out
- **SAGA Coreografada com Compensações** — Cancelamentos disparam transações compensatórias (ex: restauração de estoque) quando aplicável, garantindo consistência eventual entre contextos

A aplicação segue os princípios de **Domain-Driven Design (DDD)** e **Arquitetura Hexagonal**, com organização por contextos de domínio independentes.

## Fluxos Principais

- **Repair Order — fluxo end-to-end**  
  Fluxo operacional completo da Ordem de Serviço, da criação à entrega do veículo, com foco na ordem de chamadas dos endpoints e nas regras de negócio.  
  [`docs/repair-order-complete-flow.md`](docs/repair-order-complete-flow.md)

- **Diagrama de sequência — fluxo de negócio**  
  Visão macro da interação entre Client, API Gateway, Lambda de autenticação e os módulos internos da aplicação ao longo do fluxo operacional da oficina.  
  [`docs/sequence-diagram.md`](docs/sequence-diagram.md)

> Payloads, schemas e exemplos de request/response estão documentados no **Swagger/OpenAPI**.  
> Os documentos em `docs/` focam no fluxo operacional e na visão arquitetural das interações.

## Documentação Técnica

- Glossário de domínio (termos, roles e status): [`docs/domain-glossary.md`](docs/domain-glossary.md)
- Arquitetura e decisões técnicas: [`docs/architecture.md`](docs/architecture.md)
- Modelo Entidade-Relacionamento (DER): [`docs/der.md`](docs/der.md)
- SAGA — compensações implementadas: [`docs/implemented-saga-compensations.md`](docs/implemented-saga-compensations.md)
- ADRs (Architecture Decision Records): [`docs/adr`](docs/adr)


## Aspectos Técnicos

A API é construída em **Go**, seguindo princípios de **DDD** e **Arquitetura Hexagonal**, organizada por contextos de domínio independentes.

### Arquitetura

O projeto adota **Arquitetura Hexagonal (Ports & Adapters)** combinada com **DDD**, organizado por contextos de domínio independentes.

Cada contexto é estruturado em:

- **Domain** (`domain/`) — Entidades, value objects e regras de negócio puras, sem dependências externas
- **Application** (`application/`) — Casos de uso que representam as operações do sistema (pontos de entrada) e ports (interfaces) para dependências externas
- **Infrastructure** (`infra/`) — Adapters concretos:
    - **Inbound**: handlers HTTP e listeners de eventos
    - **Outbound**: repositórios (Bun/Postgres), publicação de eventos, JWT e integrações externas

Os casos de uso são invocados por adapters de entrada e acessam recursos externos exclusivamente via **ports**, garantindo baixo acoplamento e inversão de dependência.

A comunicação entre contextos ocorre por **eventos de domínio**, publicados e consumidos através de **AWS SQS** e **AWS SNS**. Em ambiente local, SQS e SNS são emulados via **LocalStack**. Os eventos são definidos em `internal/shared/events/` e as subscrições são registradas nos respectivos `internal/bootstrap/<context>/setup.go`. Ver detalhes em [Arquitetura de Eventos (SQS + SNS)](#arquitetura-de-eventos-sqs--sns).

Para detalhes completos da arquitetura e decisões técnicas, consulte: [`docs/architecture.md`](docs/architecture.md)

### Stack Tecnológica

- **Linguagem**: Go 1.25.1
- **Framework HTTP**: Fiber v2
- **ORM**: Bun
- **Migrações**: sql-migrate
- **Testes**: Go testing + testify + godog (BDD/Gherkin)
- **Mensageria**: AWS SQS + AWS SNS (LocalStack para desenvolvimento local)
- **Containerização**: Docker + Docker Compose
- **Banco de Dados**: PostgreSQL — Escolhido por ser um SGBD relacional maduro e confiável, adequado para garantir integridade transacional em operações críticas como criação de Ordens de Serviço, aprovação de orçamentos e controle de estoque. O modelo relacional facilita a consistência entre entidades fortemente relacionadas e oferece suporte nativo a transações ACID, constraints e índices, essenciais para a confiabilidade e evolução do sistema.

## Requisitos

- Docker e Docker Compose (LocalStack e PostgreSQL sobem como containers)
- make (opcional, recomendado)

## Pipeline

O projeto utiliza **GitHub Actions** para CI/CD automatizado com os seguintes jobs:

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                      GitHub Actions - Aplicação (oficina)                                   │
│                                                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐               │
│  │   SonarCloud │──▶ │     Build    │───▶│     Push     │──▶ │   Deploy     │               │
│  │     Scan     │    │    Docker    │    │     ECR      │    │   K8s/EKS    │               │
│  └──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘               │
│   • Security           • Tests              • Tag Image         • Apply Manifests           │
│   • Code Quality       • Coverage           • Amazon ECR        • Rolling Update            │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```
1. **tests-and-quality** — Testes unitários/integração + BDD + SonarCloud
2. **docker-build** — Build da imagem Docker
3. **deploy-ecr** — Push para Amazon ECR (apenas branch `main`)
4. **deploy-k8s** — Deploy automático no Kubernetes (EKS) aplicando os manifests localizados em `deploy/k8s/` (apenas branch `main`)
   - Obtém outputs do Terraform do estado armazenado no S3
   - Configura kubectl para conectar ao cluster EKS
   - Aplica todos os recursos Kubernetes necessários
   - Executa rolling update do deployment

## Executando o Projeto

O `docker compose up` sobe PostgreSQL, LocalStack (emulação de SQS + SNS) e o container de desenvolvimento Go. Filas SQS, tópicos SNS e respectivas subscrições (fan-out) são criados automaticamente pelo script `scripts/aws/init/ready.d/01-init.sh` durante a inicialização do LocalStack. As migrações de banco são aplicadas automaticamente ao iniciar a API.

### Com Make

1. Subir os containers:
    - `make install`
    - Caso não exista `.env`, ele será criado a partir do `.env-example`
2. Aplicar migrações (opcional — a API aplica automaticamente ao iniciar):
    - `make migrate-up`
3. Rodar a API:
    - `make run`

### Sem Make

1. Criar o `.env` a partir do `.env-example`
2. Subir os containers:
    - `docker compose up -d`
    - Aguardar os containers ficarem healthy (PostgreSQL + LocalStack)
3. Aplicar migrações (opcional — a API aplica automaticamente ao iniciar):
    - `docker compose exec -T app-dev sh -lc 'test -x /go/bin/sql-migrate || GOBIN=/go/bin /usr/local/go/bin/go install github.com/rubenv/sql-migrate/sql-migrate@latest'`
    - `docker compose exec -T app-dev sh -lc '/go/bin/sql-migrate up -config=./scripts/db/dbconfig.yml -env=development'`
4. Rodar a API:
    - `docker compose exec app-dev go run ./cmd/api`

## Swagger / OpenAPI

- UI: http://localhost:8080/docs
- Spec: http://localhost:8080/openapi.yaml
- Arquivo no repositório: `assets/docs/openapi.yaml`

## Testes

- `make test` — executa testes unitários e de integração
- `make test-bdd` — executa testes BDD (Behavior-Driven Development)
- `make test-coverage` — gera relatório de cobertura
- `make sonar` — executa análise no SonarQube

Os testes de integração criam bancos isolados e aplicam migrações automaticamente.

### Last Sonar Overview:
![alt text](assets\sonar-overview.png)

### Testes BDD

O projeto inclui testes BDD escritos com **godog** (Cucumber/Gherkin) em `tests/bdd/`. Os cenários descrevem fluxos de negócio em linguagem natural (português) e validam o comportamento end-to-end dos casos de uso.

Exemplo: o fluxo de diagnóstico da Ordem de Reparo é validado desde a criação da OS até a geração automática do orçamento, incluindo cenários de erro como estoque insuficiente e transições de status inválidas.

## Deploy no Kubernetes

A aplicação é deployada automaticamente no **Amazon EKS** através do pipeline CI/CD. Os manifests Kubernetes estão localizados em `deploy/k8s/`:

### Manifests Disponíveis

- **namespace.yaml** — Cria o namespace `fiap` para isolamento dos recursos
- **serviceaccount.yaml** — Service account para o pod da aplicação
- **configmap.yaml** — Configurações não sensíveis (PORT, AWS_REGION, SQS_BASE_URL, variáveis Datadog)
- **secrets.yaml** — Dados sensíveis (JWT_SECRET, credenciais do PostgreSQL, PG_DSN)
- **app.yaml** — Deployment da aplicação com 2 réplicas iniciais, probes de health check e integração com Datadog APM
- **services.yaml** — Service do tipo NodePort expondo a aplicação externamente
- **hpa.yaml** — Horizontal Pod Autoscaler configurado para escalar de 1 a 10 réplicas baseado em CPU (target: 50%)
- **metric-server.yaml** — Metrics Server necessário para o HPA funcionar corretamente

## Arquitetura de Eventos (SQS + SNS)

A comunicação entre os bounded contexts (`repairorder`, `estimate`, `product`) e o serviço externo `payment` é totalmente orientada a eventos. Os structs e os nomes de tópico estão centralizados em `internal/shared/events/`.

### Padrão de uso

- **Filas SQS** são usadas quando o evento tem **um único consumidor** (entrega ponto a ponto). Ex.: `RepairOrderCanceled`, `EstimateRejected`, `PaymentRequest`.
- **Tópicos SNS com fan-out para SQS** são usados quando o **mesmo evento precisa ser entregue a múltiplos contextos**. O publicador envia ao tópico SNS; cada consumidor lê de uma fila SQS dedicada.

### Eventos com fan-out via SNS

| Tópico SNS | Filas SQS consumidoras |
|---|---|
| `estimate-approved` | `repairorder-estimate-approved`, `product-estimate-approved` |
| `product-stock-reduction-confirmed` | `repairorder-product-stock-reduction-confirmed`, `estimate-product-stock-reduction-confirmed` |

### Filas SQS principais

| Fila | Produtor | Consumidor |
|---|---|---|
| `repairorder-diagnostics-finished` | `repairorder` | `estimate` |
| `estimate-created` | `estimate` | `repairorder` |
| `estimate-rejected` | `estimate` | `repairorder` |
| `estimate-canceled` | `estimate` | `product` |
| `repairorder-canceled` | `repairorder` | `estimate` |
| `product-stock-insufficient-detected` | `product` | `repairorder` |
| `payment-request` | `repairorder` | `payment` (externo) |
| `payment-status-changed.fifo` | `payment` (externo) | `repairorder` |

> Em ambiente local, todas as filas, tópicos e subscrições são provisionados pelo script `scripts/aws/init/ready.d/01-init.sh`, executado automaticamente pelo LocalStack. O healthcheck do container espera o arquivo `/tmp/sqs-init.done` antes de marcar o serviço como saudável.

## SAGA Pattern e Compensações

A integração entre contextos segue o padrão **SAGA Coreografada**: não existe um orquestrador central. Cada serviço publica eventos sobre fatos do seu domínio e reage aos eventos publicados por outros. Quando uma falha ou cancelamento ocorre, **transações compensatórias** desfazem o efeito das etapas anteriores (ex.: restaurar estoque que foi reduzido na aprovação de um orçamento posteriormente cancelado).

As compensações implementadas e as lacunas mapeadas estão documentadas em:

- [`docs/implemented-saga-compensations.md`](docs/implemented-saga-compensations.md) 

## Observabilidade

A aplicação é instrumentada com **Datadog APM** para observabilidade. A integração inclui:

- **Traces** — Rastreamento distribuído de requisições HTTP via middleware do Fiber
- **Métricas customizadas** — Métricas de negócio enviadas via DogStatsD (ex: contagem de orçamentos, OS criadas)
- **Middleware HTTP** — Captura automática de latência, status code e rotas de cada request
- **Profiling** — Coleta contínua de dados de performance da aplicação Go
