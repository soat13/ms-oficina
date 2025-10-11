package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	sharederrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/maps"
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
	products, err := c.productCatalogReader.GetByIDs(ctx, maps.Keys(input.Products))
	if err != nil {
		return err
	}
	if err := c.addItemsToEstimate(estimate, input.Products, products, domain.ProductItemType); err != nil {
		return err
	}

	services, err := c.serviceCatalogReader.GetByIDs(ctx, maps.Keys(input.Services))
	if err != nil {
		return err
	}

	if err := c.addItemsToEstimate(estimate, input.Services, services, domain.ServiceItemType); err != nil {
		return err
	}

	return nil
}

func (c *Create) addItemsToEstimate(
	estimate *domain.Estimate,
	itemQuantity map[uuid.UUID]int,
	items []CatalogItemView,
	itemType domain.ItemType,
) error {

	for _, item := range items {
		err := estimate.AddItem(item.ID, item.Name, item.Price, itemQuantity[item.ID], itemType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Create) getRepairOrder(ctx context.Context, repairID uuid.UUID) (*RepairOrderView, error) {
	repairOrder, err := c.repairOrderReader.GetByID(ctx, repairID)
	if err != nil {
		return nil, err
	}

	if repairOrder == nil {
		return nil, repairorder.ErrRepairOrderNotFound
	}

	if repairOrder.Status != repairorder.StatusInDiagnostics {
		return nil, sharederrors.ErrInvalidStatusTransaction
	}

	return repairOrder, nil
}
