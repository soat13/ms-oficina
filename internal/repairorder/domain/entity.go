package domain

import (
	"time"

	"github.com/google/uuid"
)

type RepairOrderStatus string

type (
	ItemType string

	RepairOrder struct {
		ID         uuid.UUID
		CustomerID uuid.UUID
		VehicleID  uuid.UUID
		Status     RepairOrderStatus
		Services   map[uuid.UUID]int64
		Products   map[uuid.UUID]int64
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
)

const (
	ItemService ItemType = "SERVICE"
	ItemProduct ItemType = "PRODUCT"

	StatusReceived         RepairOrderStatus = "Received"
	StatusInDiagnosis      RepairOrderStatus = "InDiagnosis"
	StatusAwaitingApproval RepairOrderStatus = "AwaitingApproval"
	StatusInProgress       RepairOrderStatus = "InProgress"
	StatusCompleted        RepairOrderStatus = "Completed"
	StatusDelivered        RepairOrderStatus = "Delivered"
)

func NewRepairOrder(customerID, vehicleID uuid.UUID) (*RepairOrder, error) {
	if customerID == uuid.Nil || vehicleID == uuid.Nil {
		return nil, ErrCustomerOrVehicleInvalid
	}

	now := time.Now()
	return &RepairOrder{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     StatusReceived,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (ro *RepairOrder) touch() { ro.UpdatedAt = time.Now() }

func (ro *RepairOrder) set(itemType ItemType, id uuid.UUID, quantity int64) error {
	if ro.Status != StatusInDiagnosis {
		return ErrOperationNotAllowed
	}

	if quantity < 0 {
		return ErrQuantityMustBeNonNegative
	}

	bucket, err := ro.bucket(itemType)
	if err != nil {
		return err
	}

	if quantity == 0 {
		delete(bucket, id)
		ro.touch()
		return nil
	}

	bucket[id] = quantity
	ro.touch()
	return nil
}

func (ro *RepairOrder) getOrNewServices() map[uuid.UUID]int64 {
	if ro.Services == nil {
		ro.Services = make(map[uuid.UUID]int64)
	}
	return ro.Services
}

func (ro *RepairOrder) getOrNewProducts() map[uuid.UUID]int64 {
	if ro.Products == nil {
		ro.Products = make(map[uuid.UUID]int64)
	}
	return ro.Products
}

func (ro *RepairOrder) bucket(itemType ItemType) (map[uuid.UUID]int64, error) {
	switch itemType {
	case ItemService:
		return ro.getOrNewServices(), nil
	case ItemProduct:
		return ro.getOrNewProducts(), nil
	default:
		return nil, ErrInvalidItemType
	}
}

func (ro *RepairOrder) SetServiceQuantity(id uuid.UUID, qty int64) error {
	return ro.set(ItemService, id, qty)
}

func (ro *RepairOrder) SetProductQuantity(id uuid.UUID, qty int64) error {
	return ro.set(ItemProduct, id, qty)
}
