package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
)

type (
	RepairOrderStatus string

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

func (r *RepairOrder) MoveToAwaitingApproval() error {
	if r.Status != repairorder.StatusInDiagnostics {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = repairorder.StatusAwaitingApproval
	r.Touch()
	return nil
}

func (r *RepairOrder) MoveToApproved() error {
	if r.Status != repairorder.StatusAwaitingApproval {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = repairorder.StatusApproved
	r.Touch()
	return nil
}

func (r *RepairOrder) StartExecution() error {
	if r.Status != repairorder.StatusApproved {
		return errors.ErrInvalidStatusTransaction
	}

	r.Status = repairorder.StatusInExecution
	r.Touch()
	return nil
}
