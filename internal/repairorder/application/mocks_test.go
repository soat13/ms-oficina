package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/oficina-utils/pkg/pagination"
)

// ---------------------------------------------------------------------------
// Repository mock
// ---------------------------------------------------------------------------

type mockRepository struct {
	getByIdFn                    func(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error)
	listFn                       func(ctx context.Context, pager pagination.Pagination) ([]*domain.RepairOrder, error)
	createFn                     func(ctx context.Context, ro *domain.RepairOrder) error
	saveCancellationFn           func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfReceivedFn             func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfInDiagnosticsFn        func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfDiagnosticsFinishedFn  func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfInAwaitingApprovalFn   func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfApprovedFn             func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfInExecutionFn          func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfFinishedFn             func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfPaymentCreatedFn       func(ctx context.Context, ro *domain.RepairOrder) error
	saveIfPaymentSucceededFn     func(ctx context.Context, ro *domain.RepairOrder) error
	getAverageExecTimeFn         func(ctx context.Context) (*float64, error)
}

func (m *mockRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
	if m.getByIdFn != nil {
		return m.getByIdFn(ctx, id)
	}
	return nil, nil
}
func (m *mockRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.RepairOrder, error) {
	if m.listFn != nil {
		return m.listFn(ctx, pager)
	}
	return nil, nil
}
func (m *mockRepository) Create(ctx context.Context, ro *domain.RepairOrder) error {
	if m.createFn != nil {
		return m.createFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveCancellation(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveCancellationFn != nil {
		return m.saveCancellationFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfReceived(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfReceivedFn != nil {
		return m.saveIfReceivedFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfInDiagnostics(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfInDiagnosticsFn != nil {
		return m.saveIfInDiagnosticsFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfDiagnosticsFinished(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfDiagnosticsFinishedFn != nil {
		return m.saveIfDiagnosticsFinishedFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfInAwaitingApproval(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfInAwaitingApprovalFn != nil {
		return m.saveIfInAwaitingApprovalFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfApproved(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfApprovedFn != nil {
		return m.saveIfApprovedFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfInExecution(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfInExecutionFn != nil {
		return m.saveIfInExecutionFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfFinished(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfFinishedFn != nil {
		return m.saveIfFinishedFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfPaymentCreated(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfPaymentCreatedFn != nil {
		return m.saveIfPaymentCreatedFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) SaveIfPaymentSucceeded(ctx context.Context, ro *domain.RepairOrder) error {
	if m.saveIfPaymentSucceededFn != nil {
		return m.saveIfPaymentSucceededFn(ctx, ro)
	}
	return nil
}
func (m *mockRepository) GetAverageExecutionTime(ctx context.Context) (*float64, error) {
	if m.getAverageExecTimeFn != nil {
		return m.getAverageExecTimeFn(ctx)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// MetricsPublisher mock
// ---------------------------------------------------------------------------

type mockMetricsPublisher struct {
	createdCalled       int
	statusChangeCalled  int
	canceledCalled      int
	lastFromStatus      string
	lastToStatus        string
	executionTime       float64
	phaseDurations      map[string]float64
	integrationErrors   int
}

func (m *mockMetricsPublisher) IncRepairOrderCreated()     { m.createdCalled++ }
func (m *mockMetricsPublisher) IncRepairOrderCanceled()    { m.canceledCalled++ }
func (m *mockMetricsPublisher) IncRepairOrderStatusChange(from, to string) {
	m.statusChangeCalled++
	m.lastFromStatus = from
	m.lastToStatus = to
}
func (m *mockMetricsPublisher) RecordRepairOrderPhaseDuration(phase string, minutes float64) {
	if m.phaseDurations == nil {
		m.phaseDurations = make(map[string]float64)
	}
	m.phaseDurations[phase] = minutes
}
func (m *mockMetricsPublisher) RecordRepairOrderExecutionTime(minutes float64) {
	m.executionTime = minutes
}
func (m *mockMetricsPublisher) IncIntegrationError(integration, operation string) {
	m.integrationErrors++
}

// ---------------------------------------------------------------------------
// EventPublisher mock
// ---------------------------------------------------------------------------

type mockEventPublisher struct {
	diagnosticsFinishedFn func(ctx context.Context, id uuid.UUID, products map[uuid.UUID]int, services map[uuid.UUID]int) error
	canceledFn            func(ctx context.Context, id uuid.UUID) error
	finishedFn            func(ctx context.Context, id uuid.UUID) error
}

func (m *mockEventPublisher) PublishRepairOrderDiagnosticsFinished(ctx context.Context, id uuid.UUID, products map[uuid.UUID]int, services map[uuid.UUID]int) error {
	if m.diagnosticsFinishedFn != nil {
		return m.diagnosticsFinishedFn(ctx, id, products, services)
	}
	return nil
}

func (m *mockEventPublisher) PublishRepairOrderCanceled(ctx context.Context, id uuid.UUID) error {
	if m.canceledFn != nil {
		return m.canceledFn(ctx, id)
	}
	return nil
}

func (m *mockEventPublisher) PublishRepairOrderFinished(ctx context.Context, id uuid.UUID) error {
	if m.finishedFn != nil {
		return m.finishedFn(ctx, id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Reader mocks
// ---------------------------------------------------------------------------

type mockVehicleReader struct {
	exists bool
	err    error
}

func (m *mockVehicleReader) Exists(_ context.Context, _ uuid.UUID) (bool, error) {
	return m.exists, m.err
}

type mockCustomerReader struct {
	exists bool
	err    error
}

func (m *mockCustomerReader) Exists(_ context.Context, _ uuid.UUID) (bool, error) {
	return m.exists, m.err
}
