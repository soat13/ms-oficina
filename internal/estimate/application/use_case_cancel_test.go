package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/estimate/application"
	"github.com/soat13/ms-oficina/internal/estimate/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type trackingEventPublisher struct {
	publishCanceledCalled bool
}

func (m *trackingEventPublisher) PublishRejected(_ context.Context, _ domain.Estimate) error {
	return nil
}

func (m *trackingEventPublisher) PublishCreated(_ context.Context, _ domain.Estimate) error {
	return nil
}

func (m *trackingEventPublisher) PublishCanceled(_ context.Context, _ domain.Estimate, _ bool) error {
	m.publishCanceledCalled = true
	return nil
}

type cancelMockRepository struct {
	estimate *domain.Estimate
}

func (m *cancelMockRepository) Save(_ context.Context, _ *domain.Estimate) error { return nil }
func (m *cancelMockRepository) UpsertItem(_ context.Context, _ uuid.UUID, _ domain.Item) error {
	return nil
}
func (m *cancelMockRepository) RemoveItem(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}
func (m *cancelMockRepository) SaveIfAwaitingStock(_ context.Context, _ *domain.Estimate) error {
	return nil
}
func (m *cancelMockRepository) SaveIfAwaitingApproval(_ context.Context, _ *domain.Estimate) error {
	return nil
}
func (m *cancelMockRepository) GetByID(_ context.Context, _ uuid.UUID) (*domain.Estimate, error) {
	return m.estimate, nil
}
func (m *cancelMockRepository) GetByRepairOrderID(_ context.Context, _ uuid.UUID) (*domain.Estimate, error) {
	return m.estimate, nil
}

func newEstimateWithStatus(t *testing.T, status domain.Status) *domain.Estimate {
	t.Helper()
	now := time.Now()
	e, err := domain.NewEstimate(uuid.New(), now, now)
	require.NoError(t, err)
	e.Status = status
	return e
}

func TestCancelEstimate_PublishCanceled(t *testing.T) {
	tests := []struct {
		name                        string
		status                      domain.Status
		expectPublishCanceledCalled bool
	}{
		{
			name:                        "approved estimate publishes EstimateCanceled",
			status:                      domain.StatusApproved,
			expectPublishCanceledCalled: true,
		},
		{
			name:                        "awaiting_stock estimate publishes EstimateCanceled",
			status:                      domain.StatusAwaitingStock,
			expectPublishCanceledCalled: true,
		},
		{
			name:                        "awaiting_approval estimate does not publish EstimateCanceled",
			status:                      domain.StatusAwaitingApproval,
			expectPublishCanceledCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := newEstimateWithStatus(t, tt.status)
			publisher := &trackingEventPublisher{}
			repo := &cancelMockRepository{estimate: estimate}

			uc := application.NewCancelEstimate(repo, publisher)
			err := uc.Execute(context.Background(), application.CancelInput{ID: estimate.ID})

			require.NoError(t, err)
			assert.Equal(t, tt.expectPublishCanceledCalled, publisher.publishCanceledCalled)
		})
	}
}

func TestCancelEstimate_AlreadyCanceledOrRejected(t *testing.T) {
	for _, status := range []domain.Status{domain.StatusCanceled, domain.StatusRejected} {
		t.Run(string(status), func(t *testing.T) {
			estimate := newEstimateWithStatus(t, status)
			publisher := &trackingEventPublisher{}
			repo := &cancelMockRepository{estimate: estimate}

			uc := application.NewCancelEstimate(repo, publisher)
			err := uc.Execute(context.Background(), application.CancelInput{ID: estimate.ID})

			require.NoError(t, err)
			assert.False(t, publisher.publishCanceledCalled)
		})
	}
}

func TestCancelEstimate_NotFound(t *testing.T) {
	publisher := &trackingEventPublisher{}
	repo := &cancelMockRepository{estimate: nil}

	uc := application.NewCancelEstimate(repo, publisher)
	err := uc.Execute(context.Background(), application.CancelInput{ID: uuid.New()})

	assert.ErrorIs(t, err, application.ErrEstimateNotFound)
}
