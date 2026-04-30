package event

import (
	"context"

	"github.com/soat13/ms-oficina/internal/estimate/application"
	"github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnStockReduceConfirmed(confirmStock *application.ConfirmStock) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.StockReductionConfirmed](msg)
		if err != nil {
			return err
		}

		return confirmStock.Execute(ctx, application.ConfirmStockInput{
			EstimateID: event.EstimateID,
		})
	}
}
