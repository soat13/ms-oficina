package application

import (
	"context"
	"time"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

type HandleDiagnosticsFinished struct {
	create Create
}

func NewHandleDiagnosticsFinished(create Create) HandleDiagnosticsFinished {
	return HandleDiagnosticsFinished{
		create: create,
	}
}

func (h *HandleDiagnosticsFinished) Execute(ctx context.Context, evt events.RepairOrderDiagnosticsFinished) error {
	input := CreateInput{
		RepairOrderID: evt.RepairOrderID,
		Products:      evt.Products,
		Services:      evt.Services,
		Now:           time.Now(),
	}

	return h.create.Execute(ctx, input)
}
