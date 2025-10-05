package listeners_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application/listeners"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/estimate/mocks"
	eventBusMock "github.com/soat13/fase-1-oficina/internal/shared/eventbus/mocks"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/stretchr/testify/require"
)

func TestOnEstimateApproved_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockRepository(ctrl)
	eventBus := eventBusMock.NewMockBus(ctrl)

	estimateID := uuid.New()
	repairOrderID := uuid.New()
	now := time.Now()

	estimate := &domain.Estimate{
		ID:            estimateID,
		RepairOrderID: repairOrderID,
		Status:        domain.StatusAwaitingStock,
		Timestamps:    entity.NewTimestamps(now, now),
	}

	repository.EXPECT().
		GetByID(gomock.Any(), estimateID).
		Return(estimate, nil)

	repository.EXPECT().
		SaveIfAwaitingStock(gomock.Any(), estimate).
		DoAndReturn(func(_ context.Context, saved *domain.Estimate) (*domain.Estimate, error) {
			require.Equal(t, domain.StatusApproved, saved.Status)
			return saved, nil
		})

	eventBus.EXPECT().
		Publish(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, topic string, event interface{}) error {
			var approvedEvent estimateEvent.Approved
			err := json.Unmarshal(event.([]byte), &approvedEvent)
			require.Nil(t, err)

			require.Equal(t, "estimate.approved", topic)
			require.Equal(t, estimateID, approvedEvent.EstimateID)
			require.Equal(t, repairOrderID, approvedEvent.RepairOrderID)
			return nil
		})

	event := estimateEvent.Approved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}
	payload, err := json.Marshal(event)
	require.NoError(t, err)

	handler := listeners.OnStockReduceConfirmed(repository, eventBus)
	err = handler(context.Background(), "product.stock.reduce.confirmed", payload)
	require.NoError(t, err)

	require.Equal(t, domain.StatusApproved, estimate.Status)
}
