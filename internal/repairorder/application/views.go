package application

import (
	"time"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type RepairOrderView struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	VehicleID  uuid.UUID
	Status     sharedRepairOrder.Status
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func toView(ro *domain.RepairOrder) RepairOrderView {
	return RepairOrderView{
		ID:         ro.ID,
		CustomerID: ro.CustomerID,
		VehicleID:  ro.VehicleID,
		Status:     ro.Status,
		CreatedAt:  ro.CreatedAt,
		UpdatedAt:  ro.UpdatedAt,
	}
}
