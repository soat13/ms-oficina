package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/domain"
	sharedRepairOrder "github.com/soat13/ms-oficina/internal/shared/repairorder"
)

type RepairOrderView struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	VehicleID  uuid.UUID
	Status     sharedRepairOrder.Status
	PaymentURL *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func toView(ro *domain.RepairOrder) RepairOrderView {
	return RepairOrderView{
		ID:         ro.ID,
		CustomerID: ro.CustomerID,
		VehicleID:  ro.VehicleID,
		Status:     ro.Status,
		PaymentURL: ro.PaymentURL,
		CreatedAt:  ro.CreatedAt,
		UpdatedAt:  ro.UpdatedAt,
	}
}
