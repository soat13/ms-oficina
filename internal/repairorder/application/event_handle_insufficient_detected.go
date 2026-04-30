package application

import (
	"context"

	"github.com/soat13/ms-oficina/internal/shared/events"
)

type HandleStockInsufficient struct {
	cancel *Cancel
}

func NewHandleStockInsufficient(cancel *Cancel) HandleStockInsufficient {
	return HandleStockInsufficient{
		cancel: cancel,
	}
}

func (h *HandleStockInsufficient) Execute(ctx context.Context, evt events.StockInsufficientDetected) error {
	return h.cancel.Execute(ctx, CancelInput{RepairOrderID: evt.RepairOrderID})
}
