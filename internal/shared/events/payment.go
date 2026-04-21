package events

import (
	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/money"
)

type PaymentStatusChanged struct {
	ID            uuid.UUID `json:"id"`
	RepairOrderID uuid.UUID `json:"external_id"`
	Status        string    `json:"status"`
	PaymentURL    *string   `json:"payment_url,omitempty"`
}

func (PaymentStatusChanged) Topic() string { return "payment-status-changed.fifo" }

type PaymentRequest struct {
	RepairOrderID uuid.UUID   `json:"id"`
	Amount        money.Money `json:"amount"`
	Description   string      `json:"description"`
}

func (PaymentRequest) Topic() string { return "payment-request" }
