package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	FinishDiagnosticsInput struct {
		RepairOrderID uuid.UUID
		Products      map[uuid.UUID]int
		Services      map[uuid.UUID]int
	}

	FinishDiagnostics struct {
		repository       Repository
		eventPublisher   EventPublisher
		productReader    ProductReader
		serviceReader    ServiceReader
		metricsPublisher MetricsPublisher
	}
)

func NewFinishDiagnostics(repository Repository, eventPublisher EventPublisher, productReader ProductReader, serviceReader ServiceReader, metricsPublisher MetricsPublisher) *FinishDiagnostics {
	return &FinishDiagnostics{
		repository:       repository,
		eventPublisher:   eventPublisher,
		productReader:    productReader,
		serviceReader:    serviceReader,
		metricsPublisher: metricsPublisher,
	}
}

func (uc *FinishDiagnostics) Execute(ctx context.Context, input FinishDiagnosticsInput) error {
	if err := uc.validateServicesExist(ctx, input.Services); err != nil {
		return err
	}

	if err := uc.validateProductsAndStock(ctx, input.Products); err != nil {
		return err
	}

	repairorder, err := uc.repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}
	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.FinishDiagnostics(); err != nil {
		return err
	}

	if err := uc.repository.SaveIfInDiagnostics(ctx, repairorder); err != nil {
		return err
	}

	if repairorder.Timestamps != nil {
		uc.metricsPublisher.RecordRepairOrderPhaseDuration("in_diagnostics", time.Since(repairorder.UpdatedAt).Minutes())
	}

	return uc.eventPublisher.PublishRepairOrderDiagnosticsFinished(
		ctx,
		input.RepairOrderID,
		input.Products,
		input.Services,
	)
}

func (uc *FinishDiagnostics) validateServicesExist(ctx context.Context, services map[uuid.UUID]int) error {
	serviceIDs := maps.Keys(services)
	allExist, err := uc.serviceReader.ExistByIds(ctx, serviceIDs)
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
	productViews, err := uc.productReader.GetByIDs(ctx, productIDs)
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
