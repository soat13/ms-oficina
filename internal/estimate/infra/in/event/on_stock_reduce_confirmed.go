package event

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnStockReduceConfirmed(confirmStock *application.ConfirmStock) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.EstimateStockReductionConfirmed](msg)
		if err != nil {
			return err
		}

		return confirmStock.Execute(ctx, application.ConfirmStockInput{
			EstimateID: event.EstimateID,
		})
	}
}
