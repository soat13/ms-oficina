package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	repairOrderEvents "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder/events"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	FinishDiagnosticsInput struct {
		RepairOrderID uuid.UUID
		Products      map[uuid.UUID]int
		Services      map[uuid.UUID]int
	}

	FinishDiagnostics struct {
		Repository    Repository
		EventBus      eventbus.Bus
		ProductReader ProductReader
		ServiceReader ServiceReader
	}
)

func NewFinishDiagnostics(repository Repository, eventBus eventbus.Bus, productReader ProductReader, serviceReader ServiceReader) *FinishDiagnostics {
	return &FinishDiagnostics{
		Repository:    repository,
		EventBus:      eventBus,
		ProductReader: productReader,
		ServiceReader: serviceReader,
	}
}

func (uc *FinishDiagnostics) Execute(ctx context.Context, input FinishDiagnosticsInput) error {
	if err := uc.validateServicesExist(ctx, input.Services); err != nil {
		return err
	}

	if err := uc.validateProductsAndStock(ctx, input.Products); err != nil {
		return err
	}

	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}
	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.FinishDiagnostics(); err != nil {
		return err
	}

	if err := uc.Repository.SaveIfInDiagnostics(ctx, repairorder); err != nil {
		return err
	}

	return uc.publishEvent(ctx, input)
}

func (uc *FinishDiagnostics) validateServicesExist(ctx context.Context, services map[uuid.UUID]int) error {
	serviceIDs := maps.Keys(services)
	allExist, err := uc.ServiceReader.ExistByIds(ctx, serviceIDs)
	if err != nil {
		return err
	}
	if !allExist {
		return ErrServiceNotFound
	}

	return nil
}

func (uc *FinishDiagnostics) validateProductsAndStock(ctx context.Context, products map[uuid.UUID]int) error {
	productIDs := maps.Keys(products)
	productViews, err := uc.ProductReader.GetByIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	if len(productViews) != len(products) {
		return ErrProductNotFound
	}

	for _, product := range productViews {
		requestedQty := products[product.ID]
		if product.Stock < requestedQty {
			return ErrInsufficientStock
		}
	}

	return nil
}

func (uc *FinishDiagnostics) publishEvent(ctx context.Context, input FinishDiagnosticsInput) error {
	event := repairOrderEvents.DiagnosticsFinished{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: input.RepairOrderID,
		Products:      input.Products,
		Services:      input.Services,
	}

	b, _ := json.Marshal(event)
	return uc.EventBus.Publish(ctx, event.Topic(), b)
}
