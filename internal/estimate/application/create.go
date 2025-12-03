package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/ports/eventbus"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/money"
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

	EstimateView struct {
		ID            uuid.UUID
		RepairOrderID uuid.UUID
		Total         money.Money
	}

	CreateEstimateOutput struct {
		Estimate EstimateView
	}

	Create struct {
		repairOrderReader    RepairOrderReader
		productCatalogReader ProductCatalogReader
		serviceCatalogReader ServiceCatalogReader
		estimateRepository   Repository
		eventbus             eventbus.Bus
	}
)

func NewCreateEstimate(
	repairOrderReader RepairOrderReader,
	productCatalogReader ProductCatalogReader,
	serviceCatalogReader ServiceCatalogReader,
	repository Repository,
	eventBus eventbus.Bus,
) *Create {
	return &Create{
		repairOrderReader:    repairOrderReader,
		productCatalogReader: productCatalogReader,
		serviceCatalogReader: serviceCatalogReader,
		estimateRepository:   repository,
		eventbus:             eventBus,
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

	return c.publishEvent(ctx, *estimate)
}

func (c *Create) publishEvent(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.Created{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}

	b, _ := json.Marshal(event)
	return c.eventbus.Publish(ctx, event.Topic(), b)
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
