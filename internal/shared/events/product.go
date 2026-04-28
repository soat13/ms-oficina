package events

import (
	"time"

	"github.com/google/uuid"
)

const TopicProductEstimateApproved = "product-estimate-approved"

type StockReductionConfirmed struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (StockReductionConfirmed) Topic() string { return "product-stock-reduction-confirmed" }

type StockInsufficientDetected struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (StockInsufficientDetected) Topic() string { return "product-stock-insufficient-detected" }
