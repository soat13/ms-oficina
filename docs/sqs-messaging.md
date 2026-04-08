# SQS Messaging

## Overview

All inter-context communication in MS-Oficina flows through AWS SQS queues. The in-memory event bus was replaced by a unified `Broker` that handles both publishing and consuming messages.

The `Broker` struct (`internal/shared/infra/messaging/broker.go`) is the single component responsible for:

- **Publishing** messages to SQS queues (implements `messaging.Publisher`)
- **Subscribing** consumers to SQS queues
- **Starting** long-poll consumer goroutines
- **Shutting down** gracefully on application exit

## Architecture

```
                    ┌──────────────────────────────────┐
                    │           Broker                  │
                    │                                    │
  Use Case ──────> │  Publish(topic, payload)           │ ──────> SQS Queue
                    │                                    │
  SQS Queue ─────> │  Consumer goroutine (long-poll)    │ ──────> Event Handler ──> Use Case
                    │                                    │
                    └──────────────────────────────────┘
```

### Topic-to-Queue Mapping

The broker maps event topics to SQS queue URLs using a naming convention:

| Event Topic                          | SQS Queue Name                          |
|--------------------------------------|-----------------------------------------|
| `repairOrder.diagnostics.finished`   | `repairorder-diagnostics-finished`      |
| `repairOrder.canceled`               | `repairorder-canceled`                  |
| `repairOrder.finished`               | `repairorder-finished`                  |
| `estimate.created`                   | `estimate-created`                      |
| `estimate.approved`                  | `estimate-approved`                     |
| `estimate.rejected`                  | `estimate-rejected`                     |
| `product.stock.reduce.confirmed`     | `product-stock-reduce-confirmed`        |
| `product.stock.insufficient.detected`| `product-stock-insufficient-detected`   |
| `payment.confirmed`                  | `payment-confirmed`                     |
| `payment.failed`                     | `payment-failed`                        |

**Rule**: `strings.ToLower(strings.ReplaceAll(topic, ".", "-"))`

The full queue URL is `SQS_BASE_URL` + queue name.

## Key Files

| File | Role |
|------|------|
| `internal/shared/messaging/bus.go` | `Publisher` and `Handler` interfaces |
| `internal/shared/messaging/publisher.go` | `Publish()` helper that serializes an event and calls `Publisher.Publish` |
| `internal/shared/infra/messaging/broker.go` | `Broker` implementation (SQS client, publisher, consumer, consumer group) |
| `internal/shared/events/*.go` | Event structs with `Topic()` method |
| `internal/bootstrap/container.go` | Wires the Broker; exposes `Publisher()`, `Subscribe()`, `StartConsumers()` |
| `internal/bootstrap/*/setup.go` | Each context registers its publishers and subscribers |
| `internal/*/infra/out/event/publisher.go` | Context-specific publisher adapters |
| `internal/*/infra/in/event/*.go` | Context-specific event handlers (unmarshal + call use case) |

## How It Works at Startup

```
cmd/api/main.go
  │
  ├─ bootstrap.BuildDefault()
  │    └─ Creates Broker via NewBroker() (connects to SQS using default credential chain)
  │
  ├─ estimate.SetupDefault(container)
  │    ├─ eventPublisher = NewEventPublisher(container.Publisher())   // publishes via Broker
  │    ├─ container.Subscribe("repairOrder.canceled", handler)       // registers consumer
  │    └─ container.Subscribe("repairOrder.diagnostics.finished", handler)
  │
  ├─ repairorder.SetupDefault(container)
  │    ├─ eventPublisher = NewEventPublisher(container.Publisher())
  │    ├─ container.Subscribe("estimate.created", handler)
  │    ├─ container.Subscribe("estimate.rejected", handler)
  │    ├─ container.Subscribe("product.stock.insufficient.detected", handler)
  │    └─ container.Subscribe("product.stock.reduce.confirmed", handler)
  │
  ├─ product.SetupDefault(container)
  │    ├─ eventPublisher = NewEventPublisher(container.Publisher())
  │    └─ container.Subscribe("estimate.approved", handler)
  │
  ├─ container.StartConsumers(ctx)     // launches 1 goroutine per subscription (long-poll)
  │
  └─ container.FiberApp.Listen(":8080") // HTTP server
```

## Consumer Behavior

Each consumer runs in its own goroutine and does:

1. **Long-poll** SQS with `WaitTimeSeconds=20`, batch of 10 messages
2. For each message: call the handler with the message body
3. **On success**: delete the message from SQS
4. **On failure**: log the error, skip deletion (message returns to queue after visibility timeout)
5. After `maxReceiveCount` (3) failed attempts, the message moves to the **DLQ**
6. On context cancellation (shutdown): exit the loop

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `AWS_REGION` | AWS region | `us-east-1` |
| `AWS_ACCESS_KEY_ID` | AWS access key (from GitHub secrets in CI) | `AKIA...` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key (from GitHub secrets in CI) | `wJal...` |
| `AWS_SESSION_TOKEN` | AWS session token (optional, for temporary credentials) | `FwoG...` |
| `AWS_ENDPOINT_URL` | Custom endpoint for LocalStack | `http://localstack:4566` |
| `SQS_BASE_URL` | Base URL prefix for queue URLs | `http://localstack:4566/000000000000/` |

The `Broker` uses the AWS SDK v2 default credential chain, which automatically picks up these variables from the environment.

## Local Development with LocalStack

LocalStack emulates SQS locally. Docker Compose already has everything configured.

### Starting the environment

```bash
make up
```

This starts:
- **PostgreSQL** on port 5432
- **LocalStack** on port 4566 (with SQS)
- **app-dev** container (Go, connected to both)

LocalStack automatically runs `scripts/localstack/init-sqs.sh` on startup, which creates all 10 queues with their DLQs and a redrive policy (`maxReceiveCount: 3`).

### Verifying queues exist

```bash
docker compose exec localstack awslocal sqs list-queues
```

### Sending a test message manually

```bash
docker compose exec localstack awslocal sqs send-message \
  --queue-url http://localhost:4566/000000000000/estimate-created \
  --message-body '{"event_id":"00000000-0000-0000-0000-000000000001","occurred_at":"2026-04-06T12:00:00Z","estimate_id":"00000000-0000-0000-0000-000000000002","repair_order_id":"00000000-0000-0000-0000-000000000003"}'
```

### Checking the DLQ

```bash
docker compose exec localstack awslocal sqs receive-message \
  --queue-url http://localhost:4566/000000000000/estimate-created-dlq
```

### Running the API

```bash
make run
```

The API starts, creates the Broker (connecting to LocalStack via `AWS_ENDPOINT_URL`), registers all consumers, and begins polling.

### Running tests

Unit tests use `NewNoopBroker()` which discards all publishes and does not connect to SQS:

```bash
make test
```

No SQS infrastructure is needed for unit tests.

## Production (AWS)

In production, the application connects to real AWS SQS. The queues are provisioned via Terraform in the `soat13/oficina-infra` repository (branch `feat/fase4`).

Required environment variables (injected via GitHub secrets / Kubernetes secrets):

```
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=<from secrets>
AWS_SECRET_ACCESS_KEY=<from secrets>
AWS_SESSION_TOKEN=<from secrets, if using temporary credentials>
SQS_BASE_URL=https://sqs.us-east-1.amazonaws.com/<account-id>/
```

`AWS_ENDPOINT_URL` should **not** be set in production (the SDK uses the default AWS endpoints).

## Event Flow Example (Happy Path)

```
1. HTTP POST /repair-orders/:id/finish-diagnostics
2. RepairOrder context publishes "repairOrder.diagnostics.finished" to SQS
3. Estimate consumer picks up the message, calls CreateEstimate use case
4. Estimate context publishes "estimate.created" to SQS
5. RepairOrder consumer picks up, moves order to "awaiting_approval"
6. Client approves estimate via HTTP
7. Estimate context publishes "estimate.approved" to SQS
8. Product consumer picks up, calls ReduceStock
9. Product context publishes "product.stock.reduce.confirmed" to SQS
10. RepairOrder consumer picks up, moves order through approve -> in_execution -> finished
11. RepairOrder context publishes "repairOrder.finished" to SQS
12. (MS-Payment consumes this in a separate microservice)
```

## DLQ and Error Handling

- Every queue has a DLQ with the suffix `-dlq`
- After 3 failed processing attempts, the message is moved to the DLQ
- Handlers should be **idempotent** (SQS guarantees at-least-once delivery)
- `SaveIfStatus` repository methods already enforce expected state, providing natural idempotency
