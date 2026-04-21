package events

import (
	"time"

	"github.com/google/uuid"
)

type EstimateCreated struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (EstimateCreated) Topic() string { return "estimate-created" }

type EstimateRejected struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (EstimateRejected) Topic() string { return "estimate-rejected" }

type EstimateApproved struct {
	EventID       uuid.UUID         `json:"event_id"`
	OccurredAt    time.Time         `json:"occurred_at"`
	EstimateID    uuid.UUID         `json:"estimate_id"`
	RepairOrderID uuid.UUID         `json:"repair_order_id"`
	Products      map[uuid.UUID]int `json:"products"`
}

func (EstimateApproved) Topic() string { return "estimate-approved" }

type EstimateCanceled struct {
	EventID       uuid.UUID         `json:"event_id"`
	OccurredAt    time.Time         `json:"occurred_at"`
	EstimateID    uuid.UUID         `json:"estimate_id"`
	RepairOrderID uuid.UUID         `json:"repair_order_id"`
	Products      map[uuid.UUID]int `json:"products"`
}

func (EstimateCanceled) Topic() string { return "estimate-canceled" }

type EstimateStockReductionConfirmed struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (EstimateStockReductionConfirmed) Topic() string {
	return "estimate-product-stock-reduction-confirmed"
}
