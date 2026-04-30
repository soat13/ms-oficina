package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/estimate/domain"
	sharedErrors "github.com/soat13/ms-oficina/internal/shared/errors"
	"github.com/soat13/ms-oficina/internal/shared/repairorder"
)

type (
	ProductLine struct {
		ProductID uuid.UUID
		Qty       int
	}

	CreateInput struct {
		RepairOrderID uuid.UUID
		Products      map[uuid.UUID]int
		Services      map[uuid.UUID]int
		Now           time.Time
	}

	Create struct {
		repairOrderReader    RepairOrderReader
		productCatalogReader ProductCatalogReader
		serviceCatalogReader ServiceCatalogReader
		estimateRepository   Repository
		eventPublisher       EventPublisher
	}
)

func NewCreateEstimate(
	repairOrderReader RepairOrderReader,
	productCatalogReader ProductCatalogReader,
	serviceCatalogReader ServiceCatalogReader,
	repository Repository,
	eventPublisher EventPublisher,
) *Create {
	return &Create{
		repairOrderReader:    repairOrderReader,
		productCatalogReader: productCatalogReader,
		serviceCatalogReader: serviceCatalogReader,
		estimateRepository:   repository,
		eventPublisher:       eventPublisher,
	}
}

func (c *Create) Execute(ctx context.Context, input CreateInput) error {
	repairOrder, err := c.getRepairOrder(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	estimate, err := domain.NewEstimate(repairOrder.ID, input.Now, input.Now)
	if err != nil {
		return err
	}

	if err := c.addItemsFromRepairOrder(ctx, estimate, input); err != nil {
		return err
	}

	if err := c.estimateRepository.Save(ctx, estimate); err != nil {
		return err
	}

	return c.eventPublisher.PublishCreated(ctx, *estimate)
}

func (c *Create) addItemsFromRepairOrder(ctx context.Context, estimate *domain.Estimate, input CreateInput) error {
	itemHelper := NewItemHelper(c.productCatalogReader, c.serviceCatalogReader)
	if err := itemHelper.ValidateAndAddProducts(ctx, estimate, input.Products); err != nil {
		return err
	}

	return itemHelper.ValidateAndAddServices(ctx, estimate, input.Services)
}

func (c *Create) getRepairOrder(ctx context.Context, repairID uuid.UUID) (*RepairOrderView, error) {
	repairOrder, err := c.repairOrderReader.GetByID(ctx, repairID)
	if err != nil {
		return nil, err
	}

	if repairOrder == nil {
		return nil, repairorder.ErrRepairOrderNotFound
	}

	if repairOrder.Status != repairorder.StatusDiagnosticsFinished {
		return nil, sharedErrors.ErrInvalidStatusTransaction
	}

	return repairOrder, nil
}
