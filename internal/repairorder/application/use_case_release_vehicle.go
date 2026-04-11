package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type ReleaseVehicleInput struct {
	RepairOrderID uuid.UUID
}

type ReleaseVehicle struct {
	Repository       Repository
	metricsPublisher MetricsPublisher
}

func NewReleaseVehicle(repository Repository, metricsPublisher MetricsPublisher) *ReleaseVehicle {
	return &ReleaseVehicle{
		Repository:       repository,
		metricsPublisher: metricsPublisher,
	}
}

func (uc *ReleaseVehicle) Execute(ctx context.Context, input ReleaseVehicleInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.ReleaseVehicle(); err != nil {
		return err
	}

	if err := uc.Repository.SaveIfPaymentSucceeded(ctx, repairorder); err != nil {
		return err
	}

	if repairorder.Timestamps != nil {
		uc.metricsPublisher.RecordRepairOrderPhaseDuration("payment_succeeded", time.Since(repairorder.UpdatedAt).Minutes())
	}

	return nil
}
