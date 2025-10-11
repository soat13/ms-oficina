package listeners_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/fase-1-oficina/internal/repairorder/mocks"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
)

func TestOnEstimateApprovedOK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockRepository(ctrl)

	repairOrderID := uuid.New()
	now := time.Now()
	timestamps := entity.NewTimestamps(now, now)

	repairOrder := &domain.RepairOrder{
		ID:         repairOrderID,
		Status:     repairorder.StatusAwaitingApproval,
		Timestamps: &timestamps,
	}

	repository.EXPECT().
		GetById(gomock.Any(), repairOrderID).
		Return(repairOrder, nil)

	repository.EXPECT().
		SaveIfInAwaitingApproval(gomock.Any(), repairOrder).
		DoAndReturn(func(_ context.Context, saved *domain.RepairOrder) (*domain.RepairOrder, error) {
			require.Equal(t, repairorder.StatusApproved, saved.Status)
			return saved, nil
		})

	event := estimate.Approved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    uuid.New(),
		RepairOrderID: repairOrderID,
	}
	payload, err := json.Marshal(event)
	require.NoError(t, err)

	handler := listeners.OnEstimateApproved(repository)
	err = handler(context.Background(), "estimate.approved", payload)
	require.NoError(t, err)

	require.Equal(t, repairorder.StatusApproved, repairOrder.Status)
}

func TestOnEstimateApproved_RepoErrorOnGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repairOrderID := uuid.New()

	repository := mocks.NewMockRepository(ctrl)
	repository.EXPECT().
		GetById(gomock.Any(), repairOrderID).
		Return(nil, errors.New("db error"))

	event := estimate.Approved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    uuid.New(),
		RepairOrderID: repairOrderID,
	}
	payload, err := json.Marshal(event)
	require.NoError(t, err)

	handler := listeners.OnEstimateApproved(repository)
	err = handler(context.Background(), "estimate.approved", payload)
	require.Error(t, err)
}

func TestOnEstimateApproved_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockRepository(ctrl)

	handler := listeners.OnEstimateApproved(repository)

	payload := []byte(`{"repair_order_id":"not-a-uuid"}`)
	err := handler(context.Background(), "estimate.approved", payload)
	require.Error(t, err)
}
