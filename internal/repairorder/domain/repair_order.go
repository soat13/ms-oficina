package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/soat13/oficina-utils/pkg/entity"
	uuidHelper "github.com/soat13/oficina-utils/pkg/utils/uuid"
)

type (
	RepairOrder struct {
		ID                   uuid.UUID
		CustomerID           uuid.UUID
		VehicleID            uuid.UUID
		Status               repairorder.Status
		ExecutionTimeMinutes *int64
		PaymentURL           *string
		*entity.Timestamps
	}
)

func NewRepairOrder(id uuid.UUID, customerID, vehicleID uuid.UUID, status *repairorder.Status, createdAt, updatedAt *time.Time) (*RepairOrder, error) {
	defaultStatus := repairorder.StatusReceived

	if customerID == uuid.Nil || vehicleID == uuid.Nil {
		return nil, ErrCustomerOrVehicleIDInvalid
	}

	var timestamp entity.Timestamps
	if createdAt != nil && updatedAt != nil {
		timestamp = entity.NewTimestamps(*createdAt, *updatedAt)
	}

	if status == nil {
		status = &defaultStatus
	}

	return &RepairOrder{
		ID:         uuidHelper.IDOrNew(id),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     *status,
		Timestamps: &timestamp,
	}, nil
}

func (r *RepairOrder) UpdatePaymentURL(paymentURL *string) {
	r.PaymentURL = paymentURL
}

func (r *RepairOrder) IsCancelled() bool {
	return r.Status == repairorder.StatusCanceled
}

func (r *RepairOrder) calculateExecutionTime() {
	executionTime := time.Since(r.UpdatedAt).Minutes()
	executionTimeMinutes := int64(executionTime)
	r.ExecutionTimeMinutes = &executionTimeMinutes
}

func (r *RepairOrder) Cancel() error {
	switch r.Status {
	case repairorder.StatusFinished,
		repairorder.StatusPaymentCreated,
		repairorder.StatusPaymentSucceeded,
		repairorder.StatusPaymentFailed,
		repairorder.StatusReleased,
		repairorder.StatusCanceled:
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = repairorder.StatusCanceled
	return nil
}

func (r *RepairOrder) StartDiagnostics() error {
	return r.moveStatus(repairorder.StatusReceived, repairorder.StatusInDiagnostics)
}

func (r *RepairOrder) FinishDiagnostics() error {
	return r.moveStatus(repairorder.StatusInDiagnostics, repairorder.StatusDiagnosticsFinished)
}

func (r *RepairOrder) MoveToAwaitingApproval() error {
	return r.moveStatus(repairorder.StatusDiagnosticsFinished, repairorder.StatusAwaitingApproval)
}

func (r *RepairOrder) Approve() error {
	return r.moveStatus(repairorder.StatusAwaitingApproval, repairorder.StatusApproved)
}

func (r *RepairOrder) StartExecution() error {
	return r.moveStatus(repairorder.StatusApproved, repairorder.StatusInExecution)
}

func (r *RepairOrder) FinishExecution() error {
	r.calculateExecutionTime()
	return r.moveStatus(repairorder.StatusInExecution, repairorder.StatusFinished)
}

func (r *RepairOrder) PaymentCreated() error {
	return r.moveStatus(repairorder.StatusFinished, repairorder.StatusPaymentCreated)
}

func (r *RepairOrder) PaymentSucceeded() error {
	return r.moveStatus(repairorder.StatusPaymentCreated, repairorder.StatusPaymentSucceeded)
}

func (r *RepairOrder) PaymentFailed() error {
	return r.moveStatus(repairorder.StatusPaymentCreated, repairorder.StatusPaymentFailed)
}

func (r *RepairOrder) ReleaseVehicle() error {
	return r.moveStatus(repairorder.StatusPaymentSucceeded, repairorder.StatusReleased)
}

func (r *RepairOrder) moveStatus(statusFrom, statusTo repairorder.Status) error {
	if r.Status != statusFrom {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = statusTo
	return nil
}
