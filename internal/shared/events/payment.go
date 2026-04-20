package events

import "github.com/google/uuid"

type PaymentStatusChanged struct {
	ID         uuid.UUID `json:"id"`
	Status     string    `json:"status"`
	PaymentURL *string   `json:"payment_url,omitempty"`
}

func (PaymentStatusChanged) Topic() string { return "payment.status.changed.fifo" }
