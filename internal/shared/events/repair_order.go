package events

import (
	"time"

	"github.com/google/uuid"
)

type RepairOrderCanceled struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (RepairOrderCanceled) Topic() string { return "repairorder-canceled" }

type RepairOrderDiagnosticsFinished struct {
	EventID       uuid.UUID         `json:"event_id"`
	OccurredAt    time.Time         `json:"occurred_at"`
	RepairOrderID uuid.UUID         `json:"repair_order_id"`
	Products      map[uuid.UUID]int `json:"products"`
	Services      map[uuid.UUID]int `json:"services"`
}

func (RepairOrderDiagnosticsFinished) Topic() string { return "repairorder-diagnostics-finished" }

type RepairOrderFinished struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (RepairOrderFinished) Topic() string { return "repairorder-finished" }

type RepairOrderStockReductionConfirmed struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (RepairOrderStockReductionConfirmed) Topic() string {
	return "repairorder-product-stock-reduction-confirmed"
}
