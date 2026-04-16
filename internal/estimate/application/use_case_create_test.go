package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/soat13/oficina-utils/pkg/money"
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
			expectedError: application.ErrProductNotAvailable,
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
			expectedError: application.ErrProductNotAvailable,
		},
		{
			name:          "Non-existent product",
			productStock:  50,
			requestedQty:  50,
			expectedError: application.ErrProductNotAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID := uuid.New()
			repairOrderID := uuid.New()

			mockRepairOrderReader := &mockRepairOrderReader{
				repairOrder: &application.RepairOrderView{
					ID:     repairOrderID,
					Status: repairorder.StatusDiagnosticsFinished,
				},
			}

			var products []application.CatalogItemView
			if tt.name != "Non-existent product" {
				products = []application.CatalogItemView{
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

			serviceID := uuid.New()
			mockServiceCatalogReader := &mockServiceCatalogReader{
				services: []application.CatalogItemView{
					{
						ID:    serviceID,
						Name:  "Test Service",
						Price: money.Money{Cents: 5000},
						Stock: 0,
					},
				},
			}

			mockRepository := &mockRepository{}
			eventPublisher := &mockEventPublisher{}

			createUseCase := application.NewCreateEstimate(
				mockRepairOrderReader,
				mockProductCatalogReader,
				mockServiceCatalogReader,
				mockRepository,
				eventPublisher,
			)

			input := application.CreateInput{
				RepairOrderID: repairOrderID,
				Products:      map[uuid.UUID]int{productID: tt.requestedQty},
				Services:      map[uuid.UUID]int{serviceID: tt.requestedQty},
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
	repairOrder *application.RepairOrderView
}

func (m *mockRepairOrderReader) GetByID(ctx context.Context, id uuid.UUID) (*application.RepairOrderView, error) {
	return m.repairOrder, nil
}

type mockProductCatalogReader struct {
	products []application.CatalogItemView
}

func (m *mockProductCatalogReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return true, nil
}

func (m *mockProductCatalogReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]application.CatalogItemView, error) {
	return m.products, nil
}

type mockServiceCatalogReader struct {
	services []application.CatalogItemView
}

func (m *mockServiceCatalogReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	for _, s := range m.services {
		if s.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockServiceCatalogReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]application.CatalogItemView, error) {
	return m.services, nil
}

type mockRepository struct{}

func (m *mockRepository) Save(ctx context.Context, estimate *domain.Estimate) error {
	return nil
}

func (m *mockRepository) UpsertItem(ctx context.Context, estimateID uuid.UUID, item domain.Item) error {
	return nil
}

func (m *mockRepository) RemoveItem(ctx context.Context, estimateID uuid.UUID, itemID uuid.UUID) error {
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

type mockEventPublisher struct{}

func (m *mockEventPublisher) PublishApproved(ctx context.Context, estimate domain.Estimate) error {
	return nil
}

func (m *mockEventPublisher) PublishRejected(ctx context.Context, estimate domain.Estimate) error {
	return nil
}

func (m *mockEventPublisher) PublishCreated(ctx context.Context, estimate domain.Estimate) error {
	return nil
}

func (m *mockEventPublisher) PublishCanceled(ctx context.Context, estimate domain.Estimate) error {
	return nil
}
