package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	ProductLine struct {
		ProductID uuid.UUID
		Qty       int
	}

	CreateEstimateInput struct {
		RepairOrderID uuid.UUID
		Now           time.Time
		Products      map[uuid.UUID]int
		Services      map[uuid.UUID]int
	}

	EstimateView struct {
		ID            uuid.UUID
		RepairOrderID uuid.UUID
		Total         money.Money
	}

	CreateEstimateOutput struct {
		Estimate EstimateView
	}

	CreateEstimateFromRepairOrder struct {
		repairOrderReader    RepairOrderReader
		productCatalogReader ProductCatalogReader
		serviceCatalogReader ServiceCatalogReader
		estimateRepository   Repository
		eventbus             eventbus.Bus
	}
)

func NewCreateEstimateFromRepairOrder(
	repairOrderReader RepairOrderReader,
	productCatalogReader ProductCatalogReader,
	serviceCatalogReader ServiceCatalogReader,
	repository Repository,
	eventBus eventbus.Bus,
) *CreateEstimateFromRepairOrder {
	return &CreateEstimateFromRepairOrder{
		repairOrderReader:    repairOrderReader,
		productCatalogReader: productCatalogReader,
		serviceCatalogReader: serviceCatalogReader,
		estimateRepository:   repository,
		eventbus:             eventBus,
	}
}

func (c *CreateEstimateFromRepairOrder) Execute(ctx context.Context, input CreateEstimateInput) (*CreateEstimateOutput, error) {
	repairOrder, err := c.validRepairOrderOrError(ctx, input.RepairOrderID)
	if err != nil {
		return nil, err
	}

	estimate, err := domain.NewEstimate(repairOrder.ID, input.Now, input.Now)
	if err != nil {
		return nil, err
	}

	if err := c.addItemsFromRepairOrder(ctx, estimate, input); err != nil {
		return nil, err
	}

	if err := c.estimateRepository.Save(ctx, estimate); err != nil {
		return nil, err
	}

	_ = c.publishEvent(ctx, *estimate) // todo: add transactional outbox pattern

	return &CreateEstimateOutput{
		Estimate: toEstimateView(estimate),
	}, nil
}

func (c *CreateEstimateFromRepairOrder) publishEvent(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.Created{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}

	b, _ := json.Marshal(event)
	if err := c.eventbus.Publish(ctx, event.Topic(), b); err != nil {
		return err
	}

	return nil
}

func (c *CreateEstimateFromRepairOrder) addItemsFromRepairOrder(ctx context.Context, estimate *domain.Estimate, input CreateEstimateInput) error {
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

func (c *CreateEstimateFromRepairOrder) addItemsToEstimate(
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

func (c *CreateEstimateFromRepairOrder) validRepairOrderOrError(ctx context.Context, repairID uuid.UUID) (*RepairOrderView, error) {
	repairOrder, err := c.repairOrderReader.GetByID(ctx, repairID)
	if err != nil {
		return nil, err
	}

	if repairOrder == nil {
		return nil, ErrRepairOrderNotFound
	}

	if repairOrder.Status != repairorder.StatusInDiagnosis {
		return nil, ErrInvalidRepairOrderStatus
	}

	return repairOrder, nil
}

func toEstimateView(estimate *domain.Estimate) EstimateView {
	return EstimateView{
		ID:            estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Total:         estimate.Total(),
	}
}
