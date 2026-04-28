package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/money"
)

type StockReductionConfirmed struct {
	EventID       uuid.UUID   `json:"event_id"`
	OccurredAt    time.Time   `json:"occurred_at"`
	EstimateID    uuid.UUID   `json:"estimate_id"`
	RepairOrderID uuid.UUID   `json:"repair_order_id"`
	TotalEstimate money.Money `json:"total_estimate"`
}

func (StockReductionConfirmed) Topic() string { return "product-stock-reduction-confirmed" }

type StockInsufficientDetected struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (StockInsufficientDetected) Topic() string { return "product-stock-insufficient-detected" }
