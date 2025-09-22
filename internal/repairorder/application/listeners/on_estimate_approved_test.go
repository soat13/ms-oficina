package listeners_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/fase-1-oficina/internal/repairorder/mocks"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
)

func TestOnEstimateApproved_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockRepository(ctrl)

	repairOrderID := uuid.New()
	now := time.Now()

	repairOrder := &domain.RepairOrder{
		ID:         repairOrderID,
		Status:     repairorder.StatusAwaitingApproval,
		Timestamps: entity.NewTimestamps(now, now),
	}

	repository.EXPECT().
		GetById(gomock.Any(), repairOrderID).
		Return(repairOrder, nil)

	repository.EXPECT().
		Save(gomock.Any(), repairOrder).
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

	handler := listeners.OnEstimateApprovedByCustomer(repository)
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

	handler := listeners.OnEstimateApprovedByCustomer(repository)
	err = handler(context.Background(), "estimate.approved", payload)
	require.Error(t, err)
}

func TestOnEstimateApproved_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mocks.NewMockRepository(ctrl)

	handler := listeners.OnEstimateApprovedByCustomer(repository)

	payload := []byte(`{"repair_order_id":"not-a-uuid"}`)
	err := handler(context.Background(), "estimate.approved", payload)
	require.Error(t, err)
}
