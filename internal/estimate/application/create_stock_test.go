package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	"github.com/soat13/fase-1-oficina/pkg/money"
	"github.com/stretchr/testify/assert"
)

func TestCreateStockValidation(t *testing.T) {
	tests := []struct {
		name          string
		productStock  int
		requestedQty  int
		expectedError error
	}{
		{
			name:          "Sufficient stock available",
			productStock:  50,
			requestedQty:  50,
			expectedError: nil,
		},
		{
			name:          "Insufficient stock",
			productStock:  50,
			requestedQty:  100,
			expectedError: ErrProductNotAvailable,
		},
		{
			name:          "Exact stock match",
			productStock:  25,
			requestedQty:  25,
			expectedError: nil,
		},
		{
			name:          "Zero stock",
			productStock:  0,
			requestedQty:  1,
			expectedError: ErrProductNotAvailable,
		},
		{
			name:          "Non-existent product",
			productStock:  50,
			requestedQty:  50,
			expectedError: ErrProductNotAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID := uuid.New()
			repairOrderID := uuid.New()

			mockRepairOrderReader := &mockRepairOrderReader{
				repairOrder: &RepairOrderView{
					ID:     repairOrderID,
					Status: "in_diagnostics",
				},
			}

			var products []CatalogItemView
			if tt.name != "Non-existent product" {
				products = []CatalogItemView{
					{
						ID:    productID,
						Name:  "Test Product",
						Price: money.Money{Cents: 1000},
						Stock: tt.productStock,
					},
				}
			}

			mockProductCatalogReader := &mockProductCatalogReader{
				products: products,
			}

			mockServiceCatalogReader := &mockServiceCatalogReader{
				services: []CatalogItemView{},
			}

			mockRepository := &mockRepository{}
			mockEventBus := &mockEventBus{}

			createUseCase := NewCreateEstimate(
				mockRepairOrderReader,
				mockProductCatalogReader,
				mockServiceCatalogReader,
				mockRepository,
				mockEventBus,
			)

			input := CreateInput{
				RepairOrderID: repairOrderID,
				Products:      map[uuid.UUID]int{productID: tt.requestedQty},
				Services:      map[uuid.UUID]int{},
				Now:           time.Now(),
			}

			err := createUseCase.Execute(context.Background(), input)

			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

type mockRepairOrderReader struct {
	repairOrder *RepairOrderView
}

func (m *mockRepairOrderReader) GetByID(ctx context.Context, id uuid.UUID) (*RepairOrderView, error) {
	return m.repairOrder, nil
}

type mockProductCatalogReader struct {
	products []CatalogItemView
}

func (m *mockProductCatalogReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error) {
	return m.products, nil
}

type mockServiceCatalogReader struct {
	services []CatalogItemView
}

func (m *mockServiceCatalogReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error) {
	return m.services, nil
}

type mockRepository struct{}

func (m *mockRepository) Save(ctx context.Context, estimate *domain.Estimate) error {
	return nil
}

func (m *mockRepository) SaveIfAwaitingStock(ctx context.Context, estimate *domain.Estimate) error {
	return nil
}

func (m *mockRepository) SaveIfAwaitingApproval(ctx context.Context, estimate *domain.Estimate) error {
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Estimate, error) {
	return nil, nil
}

func (m *mockRepository) GetByRepairOrderID(ctx context.Context, repairOrderID uuid.UUID) (*domain.Estimate, error) {
	return nil, nil
}

type mockEventBus struct{}

func (m *mockEventBus) Publish(ctx context.Context, topic string, payload []byte) error {
	return nil
}

func (m *mockEventBus) Subscribe(topic string, h eventbus.Handler) {
}
