package events

import (
	"time"

	"github.com/google/uuid"
)

type StockReduceConfirmed struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (StockReduceConfirmed) Topic() string { return "product.stock.reduce.confirmed" }

type StockInsufficientDetected struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (StockInsufficientDetected) Topic() string { return "product.stock.insufficient.detected" }
