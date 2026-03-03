DER - Modelo de Dados

Este diagrama representa a estrutura final do banco PostgreSQL, considerando todas as migrations aplicadas.

A modelagem reflete o fluxo completo da oficina:
```
Cliente → Veículo
Cliente/Veículo → Ordem de Serviço
Ordem de Serviço → Orçamento
Orçamento → Itens (Serviços ou Produtos)
```

> **Nota:** Este diagrama está em formato Mermaid. Dependendo da plataforma, ele pode ser exibido como diagrama visual ou apenas como bloco de código.
```mermaid
erDiagram
  CUSTOMERS ||--o{ VEHICLES : has
  CUSTOMERS ||--o{ REPAIR_ORDERS : opens
  VEHICLES  ||--o{ REPAIR_ORDERS : used_in
  REPAIR_ORDERS ||--o{ ESTIMATES : generates
  ESTIMATES ||--o{ ESTIMATE_ITEMS : contains

  CUSTOMERS {
    uuid id PK
    varchar name
    varchar document UK
    varchar document_type
    varchar email UK
    varchar phone_number
    timestamp created_at
    timestamp updated_at
  }

  VEHICLES {
    uuid id PK
    uuid customer_id FK
    varchar plate UK
    varchar brand
    varchar model
    int year
    timestamp created_at
    timestamp updated_at
  }

  SERVICES {
    uuid id PK
    varchar name
    bigint price
    varchar currency
    timestamp created_at
    timestamp updated_at
  }

  PRODUCTS {
    uuid id PK
    varchar name
    bigint price
    int stock
    timestamp created_at
    timestamp updated_at
  }

  REPAIR_ORDERS {
    uuid id PK
    uuid customer_id FK
    uuid vehicle_id FK
    varchar status
    bigint execution_time_minutes
    timestamp created_at
    timestamp updated_at
  }

  ESTIMATES {
    uuid id PK
    uuid repair_order_id FK
    varchar status
    timestamp created_at
    timestamp updated_at
  }

  ESTIMATE_ITEMS {
    uuid id PK
    uuid estimate_id FK
    uuid item_id
    varchar item_name
    varchar item_type
    bigint price
    int quantity
    timestamp created_at
    timestamp updated_at
  }

  USERS {
    uuid id PK
    varchar name
    varchar document UK
    varchar document_type
    varchar email UK
    varchar phone_number
    varchar password
    user_role_array roles
    timestamp created_at
    timestamp updated_at
  }
```