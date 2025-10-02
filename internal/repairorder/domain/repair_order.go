package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
)

type (
	RepairOrder struct {
		ID         uuid.UUID
		CustomerID uuid.UUID
		VehicleID  uuid.UUID
		Status     repairorder.Status
		entity.Timestamps
	}
)

func NewRepairOrder(customerID, vehicleID uuid.UUID) (*RepairOrder, error) {
	if customerID == uuid.Nil || vehicleID == uuid.Nil {
		return nil, ErrCustomerOrVehicleIDInvalid
	}

	now := time.Now()

	return &RepairOrder{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     repairorder.StatusReceived,
		Timestamps: entity.NewTimestamps(now, now),
	}, nil
}
func (r *RepairOrder) FinishExecution() error {
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
	r.Touch()
	return nil
}
