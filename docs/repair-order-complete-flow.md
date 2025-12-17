# Fluxo End-to-End — Repair Order (Ordem de Reparo)

Este documento descreve **a ordem correta de chamadas dos endpoints** para executar o ciclo completo de uma Ordem de Reparo (OS), do cadastro inicial até a entrega do veículo.

> **Importante:** Os detalhes de payloads, schemas e exemplos de request/response estão documentados no **Swagger/OpenAPI** do projeto. Este documento foca **exclusivamente no fluxo operacional** e nas regras de negócio envolvidas.

---

## Pré-requisitos

1. Autenticar via `POST /auth/login`
2. Usar o token retornado no header:

Authorization: Bearer <token>


Todas as rotas abaixo pertencem ao escopo `/admin/**`.

---

## 1. Preparação do ambiente (cadastros obrigatórios)

Antes de criar uma Ordem de Reparo, os seguintes cadastros **devem existir**:

### 1.1 Cliente
- Criar cliente:  
  `POST /admin/customers`
- Listar/buscar clientes:  
  `GET /admin/customers`  
  `GET /admin/customers/{id}`

### 1.2 Veículo (vinculado ao cliente)
- Criar veículo:  
  `POST /admin/vehicles` (informando `customer_id`)
- Listar veículos de um cliente:  
  `GET /admin/customers/{customerId}/vehicles`

### 1.3 Catálogo
- Serviços:  
  `POST /admin/services`  
  `GET /admin/services`
- Produtos e estoque:  
  `POST /admin/products`  
  `PATCH /admin/products/{id}` (ajuste de estoque)

> Sem produtos, serviços ou estoque suficiente, o fluxo poderá falhar na etapa de diagnóstico.

---

## 2. Criação da Ordem de Reparo (OS)

- Criar OS:  
  `POST /admin/repair-orders`
- Consultar OS:  
  `GET /admin/repair-orders/{id}`

**Resultado esperado:**  
A OS é criada com status inicial **Received**.

---

## 3. Diagnóstico do veículo

### 3.1 Iniciar diagnóstico
- `POST /admin/repair-orders/{id}/start-diagnostics`

Utilize este endpoint quando o mecânico iniciar a avaliação do veículo.

**Efeito esperado:**  
A OS passa para o status **In Diagnostics**.

---

### 3.2 Finalizar diagnóstico
- `POST /admin/repair-orders/{id}/finish-diagnostics`

Nesta etapa são informados:
- Produtos necessários
- Serviços necessários

**Efeitos esperados:**
- A OS muda para **Diagnostics Finished / Awaiting Approval**
- O orçamento (Estimate) é criado ou atualizado com base nos itens informados

---

## 4. Decisão do cliente (Orçamento)

Após o diagnóstico, o cliente decide sobre o orçamento.

### 4.1 Aprovar orçamento
- `POST /admin/repair-orders/{id}/estimate/approve`

**Efeitos esperados:**
- Orçamento passa para **Approved**
- OS passa para **Approved**
- Validação e/ou baixa de estoque ocorre nesta etapa

---

### 4.2 Rejeitar orçamento
- `POST /admin/repair-orders/{id}/estimate/reject`

**Efeitos esperados:**
- Orçamento passa para **Rejected**
- OS é encerrada ou cancelada, sem seguir para execução

---

## 5. Execução do serviço

### 5.1 Iniciar execução
- `POST /admin/repair-orders/{id}/start-execution`

Utilize quando os serviços aprovados começarem a ser executados.

**Efeito esperado:**  
OS passa para **In Execution**.

---

### 5.2 Finalizar execução
- `POST /admin/repair-orders/{id}/finish-execution`

Utilize quando todos os serviços forem concluídos.

**Efeito esperado:**  
OS passa para **Finished**.

---

## 6. Liberação do veículo

- `POST /admin/repair-orders/{id}/release-vehicle`

Utilize quando o veículo for entregue ao cliente.

**Efeito esperado:**  
OS passa para **Released** (estado final).

---

## 7. Cancelamento (fluxo alternativo)

- `POST /admin/repair-orders/{id}/cancel`

Pode ser utilizado quando:
- O cliente desiste do serviço
- O orçamento é rejeitado
- Alguma regra de negócio impede a continuidade do fluxo

**Efeito esperado:**  
OS passa para **Canceled**.

---

## 8. Acompanhamento e métricas

- Listar ordens de reparo:  
  `GET /admin/repair-orders`
- Consultar tempo médio de execução:  
  `GET /admin/repair-orders/average-execution-time`

Esses endpoints auxiliam no acompanhamento operacional e análise de performance da oficina.

---

## Resumo do fluxo (ordem de chamadas)

1. Criar cliente e veículo
2. Criar produtos e serviços
3. Criar Ordem de Reparo
4. Iniciar diagnóstico
5. Finalizar diagnóstico (gera orçamento)
6. Aprovar **ou** rejeitar orçamento
7. (Se aprovado) iniciar execução
8. Finalizar execução
9. Liberar veículo

---

Este fluxo representa o **happy path completo** da Ordem de Reparo na API da Oficina.
