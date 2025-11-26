package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	AddItemInput struct {
		EstimateID uuid.UUID
		ItemID     uuid.UUID
		Quantity   int
	}

	AddItem struct {
		productCatalogReader ProductCatalogReader
		serviceCatalogReader ServiceCatalogReader
		repository           Repository
	}
)

func NewAddItem(
	productCatalogReader ProductCatalogReader,
	serviceCatalogReader ServiceCatalogReader,
	repository Repository,
) *AddItem {
	return &AddItem{
		productCatalogReader: productCatalogReader,
		serviceCatalogReader: serviceCatalogReader,
		repository:           repository,
	}
}

func (a *AddItem) Execute(ctx context.Context, input AddItemInput) error {
	estimate, err := a.repository.GetByID(ctx, input.EstimateID)
	if err != nil {
		return err
	}
	if estimate == nil {
		return ErrEstimateNotFound
	}

	itemMap := map[uuid.UUID]int{input.ItemID: input.Quantity}
	existProduct, err := a.productCatalogReader.Exists(ctx, input.ItemID)
	if err != nil {
		return err
	}
	existService, err := a.serviceCatalogReader.Exists(ctx, input.ItemID)
	if err != nil {
		return err
	}

	if !existProduct && !existService {
		return ErrProductOrServiceNotFound
	}

	itemHelper := NewItemHelper(a.productCatalogReader, a.serviceCatalogReader)
	if existProduct {
		if err := itemHelper.ValidateAndAddProducts(ctx, estimate, itemMap); err != nil {
			return err
		}
	} else {
		if err := itemHelper.ValidateAndAddServices(ctx, estimate, itemMap); err != nil {
			return err
		}
	}

	return a.repository.Save(ctx, estimate)
}
