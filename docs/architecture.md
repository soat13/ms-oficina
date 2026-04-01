# Arquitetura

Este projeto segue **Arquitetura Hexagonal (Ports & Adapters)** combinada com **Domain-Driven Design (DDD)**, com foco no isolamento das regras de negócio, na redução de acoplamento e na evolução incremental do sistema.

A aplicação é organizada por **contextos de domínio independentes**, cada um representando um bounded context do negócio.

---

## Visão Geral

A arquitetura separa claramente:

- regras de negócio
- casos de uso (orquestração)
- detalhes de infraestrutura

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                             INBOUND (ADAPTERS)                               │
│                                                                              │
│  Entradas do sistema                                                         │
│                                                                              │
│  • HTTP Handlers                                                             │
│    internal/*/infra/in/http/handler.go                                       │
│                                                                              │
│  • Event Listeners                                                           │
│    internal/*/infra/in/event/on_*.go                                         │
└───────────────────────────────────────┬──────────────────────────────────────┘
                                        │ calls
                                        ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                         HEXAGON CORE (APPLICATION)                           │
│                                                                              │
│  Inbound side (operations / entry points)                                    │
│                                                                              │
│  • Use Cases                                                                 │
│    internal/*/application/use_case_*.go                                      │
│                                                                              │
│  • Event Handlers                                                            │
│    internal/*/application/event_handle_*.go                                  │
│                                                                              │
│  Domain Model                                                                │
│                                                                              │
│  • Entities / Value Objects / Business Rules                                 │
│    internal/*/domain/*                                                       │
│                                                                              │
│  Outbound side (ports = interfaces)                                          │
│                                                                              │
│  • Application Ports                                                         │
│    internal/*/application/ports.go                                           │
│                                                                              │
│  • Shared Event Port                                                         │
│    internal/ports/event/*                                                    │
└───────────────────────────────────────┬──────────────────────────────────────┘
                                        │ depends on (via outbound ports)
                                        ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                             OUTBOUND (ADAPTERS)                              │
│                                                                              │
│  Saídas do sistema                                                           │
│                                                                              │
│  • Database Repositories / Readers                                           │
│    internal/*/infra/out/db/*                                                 │
│                                                                              │
│  • Event Publisher                                                           │
│    internal/*/infra/out/event/publisher.go                                   │
│                                                                              │
│  • Token / Authentication Adapter                                            │
│    internal/auth/infra/out/jwt/*                                             │
└──────────────────────────────────────────────────────────────────────────────┘

```
O diagrama acima representa a aplicação da Arquitetura Hexagonal no projeto,
enquanto as seções a seguir detalham as responsabilidades e decisões de cada camada.

---

## Estrutura por Camada

### Domain
`internal/<contexto>/domain`

Responsável por:
- Entidades
- Value Objects
- Regras de negócio
- Erros de domínio

Características:
- Não depende de `application` nem `infra`
- Não conhece frameworks, banco de dados ou protocolos
- Contém apenas lógica de negócio pura

---

### Application (Use Cases)
`internal/<contexto>/application`

Responsável por:
- Implementar os casos de uso do sistema, responsáveis por orquestrar o fluxo das regras de negócio
- Coordenar interações entre entidades
- Definir **ports** (interfaces) para acesso a recursos externos

Cada **use case**:
- representa uma operação específica do sistema
- é um ponto de entrada independente
- possui uma única responsabilidade

> O projeto **não utiliza um service agregador por contexto**.  
> Essa decisão evita acoplamento indevido entre casos de uso e respeita o princípio de responsabilidade única (SRP) e segregação de interfaces (ISP).

---

### Ports (Outbound)
`internal/<contexto>/application/ports.go`  
`internal/ports/*`

Os ports são interfaces que representam dependências técnicas necessárias aos casos de uso, como:

- repositórios
- leitores
- publicação de eventos
- geração/validação de tokens

Essas interfaces pertencem à camada `application` e permitem que os casos de uso permaneçam independentes de detalhes de infraestrutura.

---

### Infrastructure (Adapters)
`internal/<contexto>/infra/in/*`  
`internal/<contexto>/infra/out/*`

A camada de infraestrutura contém adapters concretos, responsáveis por integrar a aplicação a mecanismos de entrada e saída.

#### Adapters de Entrada (Inbound)

- **HTTP**: `internal/<contexto>/infra/in/http`
    - Recebe requisições via Fiber
    - Realiza parsing e validação de dados
    - Invoca os use cases da camada `application`

- **Eventos (consumo)**: `internal/<contexto>/infra/in/event/on_*.go`
    - Escuta eventos publicados por outros contextos
    - Converte o evento em uma chamada para a camada `application`

A lógica de reação a eventos (regras e orquestração) permanece na camada `application`, garantindo que a infraestrutura não contenha regras de negócio.

---

#### Adapters de Saída (Outbound)

- **Persistência**: `internal/<contexto>/infra/out/db`
    - Implementa repositórios e leitores usando Bun/Postgres

- **Eventos (publicação)**: `internal/<contexto>/infra/out/event/publisher.go`
    - Implementa o port de publicação de eventos
    - Utiliza o event bus via pkg externo

- **Autenticação (JWT)**: `internal/auth/infra/out/jwt`
    - Implementa geração e validação de tokens

Esses adapters implementam os ports definidos na camada `application`.

---

## Bootstrap e Composição
`internal/bootstrap`

Responsável por:
- Criar instâncias dos use cases
- Injetar dependências (ports e adapters)
- Registrar handlers HTTP e listeners de eventos
- Inicializar o sistema

Toda a composição acontece fora das camadas de domínio e aplicação, mantendo as regras de negócio desacopladas de detalhes técnicos.

---

## Fluxos Principais

### Fluxo HTTP (exemplo)

1. Handler HTTP recebe a requisição
2. Converte input para o formato esperado pelo use case
3. Invoca o use case na camada `application`
4. Use case executa a regra de negócio
5. Acesso a recursos externos ocorre via ports
6. Resposta é traduzida para HTTP pelo handler

---

### Fluxo de Eventos

1. Um contexto publica um evento de domínio
2. O evento é entregue pelo event bus
3. Um listener (`infra/event/on_*.go`) recebe o evento
4. O listener atua apenas como adapter de entrada; a lógica de reação ao evento permanece na camada `application`.
5. Novos eventos podem ser publicados via ports

---

## Comunicação entre Contextos

A comunicação entre bounded contexts ocorre exclusivamente por **eventos de domínio**, evitando dependências diretas entre contextos e preservando o desacoplamento arquitetural.

Atualmente, o projeto utiliza um **event bus em memória**, permitindo fácil substituição futura por uma implementação externa (ex: Kafka, RabbitMQ) sem impacto na camada de aplicação.

---

## Considerações Finais

Essa abordagem garante:
- isolamento das regras de negócio
- facilidade de teste
- baixo acoplamento
- evolução incremental da arquitetura

A Arquitetura Hexagonal foi adotada por favorecer o isolamento das regras de negócio, facilitar testes automatizados e permitir a substituição de detalhes técnicos sem impacto na camada de aplicação.
