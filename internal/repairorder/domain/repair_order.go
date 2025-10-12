package domain

import (
	"time"

	"github.com/google/uuid"
	uuidHelper "github.com/soat13/fase-1-oficina/pkg/utils/uuid"

	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
)

type (
	RepairOrder struct {
		ID                   uuid.UUID
		CustomerID           uuid.UUID
		VehicleID            uuid.UUID
		Status               repairorder.Status
		ExecutionTimeMinutes *int64
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

func (r *RepairOrder) Cancel() error {
	if r.Status == repairorder.StatusReleased {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = repairorder.StatusCanceled
	return nil
}

func (r *RepairOrder) IsCancelled() bool {
	return r.Status == repairorder.StatusCanceled
}

func (r *RepairOrder) FinishExecution() error {
	r.calculateExecutionTime()
	return r.moveStatus(repairorder.StatusInExecution, repairorder.StatusFinished)
}

func (r *RepairOrder) MoveToAwaitingApproval() error {
	return r.moveStatus(repairorder.StatusInDiagnostics, repairorder.StatusAwaitingApproval)
}

func (r *RepairOrder) MoveToApproved() error {
	return r.moveStatus(repairorder.StatusAwaitingApproval, repairorder.StatusApproved)
}

func (r *RepairOrder) StartExecution() error {
	return r.moveStatus(repairorder.StatusApproved, repairorder.StatusInExecution)
}

func (r *RepairOrder) ReleaseVehicle() error {
	return r.moveStatus(repairorder.StatusFinished, repairorder.StatusReleased)
}

func (r *RepairOrder) moveStatus(statusFrom, statusTo repairorder.Status) error {
	if r.Status != statusFrom {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = statusTo
	return nil
}

func (r *RepairOrder) calculateExecutionTime() {
	executionTime := time.Since(r.UpdatedAt).Minutes()
	executionTimeMinutes := int64(executionTime)
	r.ExecutionTimeMinutes = &executionTimeMinutes
}
