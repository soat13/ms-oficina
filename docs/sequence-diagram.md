Diagrama de Sequência - Fluxo de Negócio

Este diagrama representa o fluxo completo de interação entre o cliente, API Gateway, Lambda de autenticação e os módulos internos da aplicação, considerando os endpoints expostos pela Oficina API.

A modelagem reflete o fluxo operacional da oficina:

```text
Autenticação → CRUDs Administrativos → Ordem de Serviço
Ordem de Serviço → Diagnóstico → Verificação de Estoque
Ordem de Serviço → Aprovação/Rejeição do Orçamento
Aprovação → Execução → Finalização → Liberação do Veículo
```

O diagrama demonstra:

- Interação entre Client, API Gateway e Lambda de Autenticação
- Chamadas para os endpoints administrativos (/admin/...)
- Transições de status da entidade RepairOrder
- Validações de estoque durante o processo de diagnóstico
- Aprovação ou rejeição do orçamento
- Processo completo até a devolução do veículo

> Nota 1: Este diagrama está em formato Mermaid. Dependendo da plataforma, ele pode ser exibido como diagrama visual ou apenas como bloco de código.

> Nota 2: O diagrama representa o fluxo conceitual da aplicação. Camadas internas como UseCases, Repositórios, Transações e Validações de Domínio podem estar abstraídas para simplificação visual.

> Nota 3: As respostas HTTP (204, 200, 201, 422 etc.) refletem o comportamento atual dos handlers implementados, podendo evoluir conforme as regras de negócio sejam ajustadas.

```mermaid
sequenceDiagram
autonumber
actor Client
participant APIGW as API Gateway
participant Auth as Auth Lambda
participant Users as Users Endpoint
participant Customers as Customers Endpoint
participant Vehicles as Vehicles Endpoint
participant Services as Services Endpoint
participant Products as Products Endpoint
participant Repair as Repair Orders Endpoint
participant Estimates as Estimates Endpoint
participant Stock as Stock Module

    %% =========================
    %% AUTHENTICATION
    %% =========================
    Client->>APIGW: POST /auth/login
    APIGW->>Auth: Invoke login
    Auth-->>APIGW: 200 {access_token}
    APIGW-->>Client: 200 {access_token}

    Note over Client: Todas próximas chamadas usam Bearer Token

    %% =========================
    %% USERS CRUD (exemplo)
    %% =========================
    Client->>APIGW: POST /admin/users
    APIGW->>Users: Create User
    Users-->>APIGW: 201
    APIGW-->>Client: 201

    %% =========================
    %% CUSTOMERS CRUD (exemplo)
    %% =========================
    Client->>APIGW: POST /admin/customers
    APIGW->>Customers: Create Customer
    Customers-->>APIGW: 201
    APIGW-->>Client: 201

    %% =========================
    %% VEHICLES CRUD (exemplo)
    %% =========================
    Client->>APIGW: POST /admin/vehicles
    APIGW->>Vehicles: Create Vehicle
    Vehicles-->>APIGW: 201
    APIGW-->>Client: 201

    %% =========================
    %% SERVICES CRUD (exemplo)
    %% =========================
    Client->>APIGW: POST /admin/services
    APIGW->>Services: Create Service
    Services-->>APIGW: 201
    APIGW-->>Client: 201

    %% =========================
    %% PRODUCTS CRUD (exemplo)
    %% =========================
    Client->>APIGW: POST /admin/products
    APIGW->>Products: Create Product
    Products-->>APIGW: 201
    APIGW-->>Client: 201

    %% ======================================================
    %% REPAIR ORDER FLOW
    %% ======================================================

    Note over Client,Repair: Criar Ordem de Serviço (Status: received)

    Client->>APIGW: POST /admin/repair-orders
    APIGW->>Repair: Create Repair Order
    Repair-->>APIGW: 201
    APIGW-->>Client: 201

    %% START DIAGNOSTIC
    Note over Client,Repair: Iniciar Diagnóstico (received -> in_diagnostics)

    Client->>APIGW: POST /admin/repair-orders/{id}/start-diagnostics
    APIGW->>Repair: start-diagnostics
    Repair-->>APIGW: 204
    APIGW-->>Client: 204

    %% FINISH DIAGNOSTIC + STOCK CHECK
    Note over Client,Repair: Finalizar Diagnóstico (in_diagnostics -> diagnostics_finished)\n+ Verifica estoque antes de liberar orçamento para aprovação

    Client->>APIGW: POST /admin/repair-orders/{id}/finish-diagnostics {products[], services[]}
    APIGW->>Repair: finish-diagnostics (valida itens)
    Repair->>Stock: CheckStock(products, quantities)
    alt estoque insuficiente / produto indisponível
        Stock-->>Repair: ErrInsufficientStock / ErrProductNotAvailable
        Repair-->>Client: Notifica usuário sobre itens com problemas de estoque
        Repair-->>APIGW: 204
        APIGW-->>Client: 204
    else estoque OK
        Stock-->>Repair: OK
        Repair-->>APIGW: 204
        APIGW-->>Client: 204
        Note over Repair: Status interno segue para awaiting_approval\n(e orçamento/estimate fica "aberto")
    end

    %% ADD ITEMS TO ESTIMATE (quando awaiting_approval)
    loop Alterar itens enquanto orçamento awaiting_approval
        Client->>APIGW: POST /admin/estimates/{id}/items
        APIGW->>Estimates: Add item
        Estimates-->>APIGW: 200
        APIGW-->>Client: 200

        Client->>APIGW: DELETE /admin/estimates/{id}/items/{itemId}
        APIGW->>Estimates: Remove item
        Estimates-->>APIGW: 200
        APIGW-->>Client: 200
    end

    %% APPROVE OR REJECT
    alt Estimate Approved
        Client->>APIGW: POST /admin/repair-orders/{id}/estimate/approve
        APIGW->>Estimates: Approve Estimate
        Estimates->>Stock: Decrease stock
        Stock-->>Estimates: OK
        Estimates-->>APIGW: 200
        APIGW-->>Client: 200

        %% START EXECUTION
        Note over Client,Repair: Iniciar Execução (approved -> in_execution)

        Client->>APIGW: POST /admin/repair-orders/{id}/start-execution
        APIGW->>Repair: start-execution
        Repair-->>APIGW: 204
        APIGW-->>Client: 204

        %% FINISH EXECUTION
        Note over Client,Repair: Finalizar Execução (in_execution -> finished)\n(calcula executionTime)

        Client->>APIGW: POST /admin/repair-orders/{id}/finish-execution
        APIGW->>Repair: finish-execution
        Repair-->>APIGW: 204
        APIGW-->>Client: 204

        %% RELEASE VEHICLE
        Note over Client,Repair: Devolver carro (finished -> released)

        Client->>APIGW: POST /admin/repair-orders/{id}/release-vehicle
        APIGW->>Repair: release-vehicle
        Repair-->>APIGW: 204
        APIGW-->>Client: 204

    else Estimate Rejected
        Client->>APIGW: POST /admin/repair-orders/{id}/estimate/reject
        APIGW->>Estimates: Reject Estimate
        Estimates-->>APIGW: 200
        APIGW-->>Client: 200

        Note over Repair: Ordem pode ser cancelada (POST /admin/repair-orders/{id}/cancel)\nou seguir regra de encerramento do domínio
    end
```