### Arquitetura

Cada contexto de domínio é organizado em três camadas:

- **Domain** (`domain/`) — Entidades, value objects, regras de negócio e erros de domínio
- **Application** (`application/`) — Casos de uso (use cases), portas (interfaces) e lógica de aplicação
- **Infrastructure** (`infra/`) — Adaptadores concretos (repositórios Bun, handlers HTTP, serviços externos)

A comunicação entre contextos ocorre através de **eventos de domínio** gerenciados por um event bus em memória, mantendo o desacoplamento entre bounded contexts.
