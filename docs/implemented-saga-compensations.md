# SAGA Pattern — Compensações Implementadas

Este documento descreve as **transações compensatórias** já implementadas no ciclo de vida da Ordem de Reparo, abrangendo os serviços `oficina` e `payment`.

---

## Visão Geral

A integração entre os serviços utiliza o padrão **SAGA Coreografado**: não há um orquestrador central; cada serviço publica eventos e reage aos eventos de outros. A consistência eventual é garantida por **transações compensatórias** — ações que desfazem o trabalho realizado por etapas anteriores quando uma falha ocorre.

**Infraestrutura de mensageria:**
- **AWS SQS**: filas para entrega ponto a ponto
- **AWS SNS**: tópicos com fan-out para múltiplos consumidores

Todos os eventos estão definidos em `internal/shared/events/`.

### Fan-out via SNS

Alguns eventos precisam ser consumidos por mais de um contexto. O SNS distribui o mesmo evento para filas SQS distintas:

| Evento publicado | Tópico SNS | Filas SQS consumidoras |
|---|---|---|
| `EstimateApproved` | `estimate-approved` | `repairorder-estimate-approved`, `product-estimate-approved` |
| `StockReductionConfirmed` | `product-stock-reduction-confirmed` | `repairorder-product-stock-reduction-confirmed`, `estimate-product-stock-reduction-confirmed` |

---

## Fluxo Feliz (Resumo)

```
RepairOrder criada (received)
  → in_diagnostics → diagnostics_finished
  → EstimateCreated → awaiting_approval
  → Estimate aprovado → EstimateApproved publicado
  → Estoque reduzido → StockReductionConfirmed
  → RepairOrder: approved → in_execution → finished
  → PaymentRequest enviado ao serviço payment
  → Pagamento processado
  → PaymentStatusChanged (SUCCEEDED) → RepairOrder: released
```

O estoque dos produtos é reduzido no momento da **aprovação do orçamento** e os produtos são fisicamente consumidos durante a execução. Por isso, compensações de estoque só fazem sentido para cancelamentos **anteriores ao início da execução**.

---

## Compensações Implementadas

### C1 — Rejeição de Orçamento

**Quando ocorre:** O gerente rejeita o orçamento via `POST /admin/estimates/{id}/reject`.

**Cadeia de eventos:**

```
[estimate context]
  Estimate rejeitado → EstimateRejected (estimate-rejected)

[repairorder context]
  HandleEstimateRejected → Cancel use case
  → RepairOrder: canceled
  → RepairOrderCanceled (repairorder-canceled)

[estimate context]
  HandleRepairOrderCanceled → CancelEstimate
  → Estimate já rejeitado → no-op idempotente
```

**Efeito líquido:**
- RepairOrder cancelada
- Estimate rejeitado
- Estoque não havia sido reduzido → nenhuma restauração necessária

**Arquivos relevantes:**
- `internal/repairorder/application/event_handle_estimate_rejected.go`
- `internal/repairorder/application/use_case_cancel.go`
- `internal/estimate/application/event_handle_repair_order_canceled.go`

---

### C2 — Estoque Insuficiente

**Quando ocorre:** A verificação automática de estoque falha logo após a aprovação do orçamento.

**Cadeia de eventos:**

```
[estimate context]
  Estimate aprovado → EstimateApproved (estimate-approved)

[product context]
  HandleEstimateApproved → ReduceStock use case
  → Estoque insuficiente para algum produto
  → StockInsufficientDetected (product-stock-insufficient-detected)

[repairorder context]
  HandleStockInsufficient → Cancel use case
  → RepairOrder: canceled
  → RepairOrderCanceled (repairorder-canceled)

[estimate context]
  HandleRepairOrderCanceled → CancelEstimate
  → Estimate estava em awaiting_stock (shouldPublishEvent = true)
  → EstimateCanceled (estimate-canceled) publicado COM Products map
  → product.RestoreStock executa (vide observação abaixo)
```

**Efeito líquido:**
- RepairOrder cancelada
- Estimate cancelado
- Nenhum produto foi descontado do estoque

**Detalhe de implementação:** O guard `shouldPublishEvent` em `use_case_cancel.go` publica `EstimateCanceled` quando o Estimate estava em `approved` ou `awaiting_stock`. Como neste cenário o Estimate está em `awaiting_stock`, o evento É emitido e `product.RestoreStock` executa (ver observação a seguir sobre o efeito em C2).

**Arquivos relevantes:**
- `internal/product/application/use_case_reduce_stock.go`
- `internal/repairorder/application/event_handle_insufficient_detected.go`
- `internal/estimate/application/use_case_cancel.go`
- `internal/estimate/application/event_handle_repair_order_canceled.go`

---

### C3 — Cancelamento Manual Após Confirmação de Estoque

**Quando ocorre:** A RepairOrder é cancelada manualmente via `POST /admin/repair-orders/{id}/cancel` após o estoque ter sido reduzido e o Estimate confirmado (`approved`), mas **antes do início da execução** (os produtos ainda não foram fisicamente consumidos).

**Cadeia de eventos:**

```
[repairorder context — HTTP]
  Cancel use case → RepairOrder: canceled
  → RepairOrderCanceled (repairorder-canceled)

[estimate context]
  HandleRepairOrderCanceled → CancelEstimate
  → Estimate estava approved (shouldPublishEvent = true)
  → EstimateCanceled (estimate-canceled) publicado COM Products map

[product context]
  HandleEstimateCanceled → RestoreStock use case
  → product.Stock += quantity (para cada produto do mapa)
  → Estoque restaurado
```

**Efeito líquido:**
- RepairOrder cancelada
- Estimate cancelado
- Produtos devolvidos ao estoque (quantidades somadas de volta)

**Detalhe de implementação:** O evento `EstimateCanceled` carrega as quantidades exatas que foram reduzidas. O `RestoreStock` use case simplesmente some essas quantidades ao estoque atual sem validação adicional — operação inversa e segura da redução original.

**Arquivos relevantes:**
- `internal/repairorder/application/use_case_cancel.go`
- `internal/estimate/application/use_case_cancel.go`
- `internal/product/application/event_handle_estimate_canceled.go`
- `internal/product/application/use_case_restore_stock.go`
- `internal/shared/events/estimate.go` — campo `EstimateCanceled.Products`

---

### C4 — Condição de Corrida Durante Redução de Estoque

**Quando ocorre:** A RepairOrder é cancelada durante a janela entre `EstimateApproved` publicado e `StockReductionConfirmed` recebido — ou seja, enquanto o Estimate está em `awaiting_stock`.

**Lacuna corrigida:** Antes desta correção, o guard só publicava `EstimateCanceled` se o Estimate estivesse em `approved`. Um Estimate em `awaiting_stock` cancelado nessa janela não emitia o evento, e o estoque reduzido fisicamente no contexto `product` nunca era restaurado.

**Cadeia de eventos:**

```
[repairorder context — HTTP]
  Cancel use case → RepairOrder: canceled
  → RepairOrderCanceled (repairorder-canceled)

[estimate context]
  HandleRepairOrderCanceled → CancelEstimate
  → Estimate estava awaiting_stock (shouldPublishEvent = true)
  → EstimateCanceled (estimate-canceled) publicado COM Products map

[product context]
  HandleEstimateCanceled → RestoreStock use case
  → product.Stock += quantity (para cada produto do mapa)
  → Estoque restaurado via UpdateBatch
```

**Efeito líquido:**
- RepairOrder cancelada
- Estimate cancelado
- Estoque restaurado, mesmo que a redução já tenha ocorrido fisicamente antes do cancelamento chegar

**Detalhe de implementação:** O guard foi alterado de `wasApproved := estimate.IsApproved()` para `shouldPublishEvent := estimate.IsApproved() || estimate.IsAwaitingStock()`. O `RestoreStock` use case é seguro para executar mesmo que a redução ainda não tenha ocorrido — se o cancelamento chegou antes da redução, a restauração prévia seria compensada pelo `StockReductionConfirmed` subsequente, que tentaria confirmar um Estimate já cancelado e seria ignorado por idempotência.

**Arquivos relevantes:**
- `internal/estimate/application/use_case_cancel.go` — guard `shouldRestoreStock`
- `internal/estimate/domain/estimate.go` — novo método `IsAwaitingStock()`
- `internal/estimate/application/use_case_cancel_test.go`
