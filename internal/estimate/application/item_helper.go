package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type itemHelper struct {
	productCatalogReader ProductCatalogReader
	serviceCatalogReader ServiceCatalogReader
}

func NewItemHelper(
	productCatalogReader ProductCatalogReader,
	serviceCatalogReader ServiceCatalogReader,
) *itemHelper {
	return &itemHelper{
		productCatalogReader: productCatalogReader,
		serviceCatalogReader: serviceCatalogReader,
	}
}

func (h *itemHelper) ValidateAndAddProducts(ctx context.Context, estimate *domain.Estimate, quantities map[uuid.UUID]int) error {
	products, err := h.productCatalogReader.GetByIDs(ctx, maps.Keys(quantities))
	if err != nil {
		return err
	}
	productMap := make(map[uuid.UUID]CatalogItemView, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	for productID, qty := range quantities {
		if qty > 0 {
			product, exists := productMap[productID]
			if !exists || product.Stock < qty {
				return ErrProductNotAvailable
			}
		}
	}

	return h.addItemsToEstimate(estimate, quantities, products, domain.ProductItemType)
}

func (h *itemHelper) ValidateAndAddServices(ctx context.Context, estimate *domain.Estimate, quantities map[uuid.UUID]int) error {
	services, err := h.serviceCatalogReader.GetByIDs(ctx, maps.Keys(quantities))
	if err != nil {
		return err
	}

	return h.addItemsToEstimate(estimate, quantities, services, domain.ServiceItemType)
}

func (h *itemHelper) addItemsToEstimate(estimate *domain.Estimate, itemQuantity map[uuid.UUID]int, items []CatalogItemView, itemType domain.ItemType) error {
	for _, item := range items {
		err := estimate.AddItem(item.ID, item.Name, item.Price, itemQuantity[item.ID], itemType)
		if err != nil {
			return err
		}
	}

	return nil
}
