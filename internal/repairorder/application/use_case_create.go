package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/domain"
)

type (
	CreateInput struct {
		CustomerID uuid.UUID
		VehicleID  uuid.UUID
	}

	CreateOutput struct {
		RepairOrderID uuid.UUID
	}

	Create struct {
		repository       Repository
		customerReader   CustomerReader
		vehicleReader    VehicleReader
		metricsPublisher MetricsPublisher
	}
)

func NewCreate(repository Repository, customerReader CustomerReader, vehicleReader VehicleReader, metricsPublisher MetricsPublisher) *Create {
	return &Create{
		repository:       repository,
		customerReader:   customerReader,
		vehicleReader:    vehicleReader,
		metricsPublisher: metricsPublisher,
	}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if err := uc.inputValidate(ctx, input); err != nil {
		return nil, err
	}

	repairOrder, err := domain.NewRepairOrder(uuid.Nil, input.CustomerID, input.VehicleID, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = uc.repository.Create(ctx, repairOrder)
	if err != nil {
		return nil, err
	}

	return &CreateOutput{RepairOrderID: repairOrder.ID}, nil
}

func (uc *Create) inputValidate(ctx context.Context, input CreateInput) error {
	exist, err := uc.vehicleReader.Exists(ctx, input.VehicleID)
	if err != nil {
		return err
	}

	if !exist {
		return ErrVehicleOrCustomerNotFound
	}

	exist, err = uc.customerReader.Exists(ctx, input.CustomerID)
	if err != nil {
		return err
	}

	if !exist {
		return ErrVehicleOrCustomerNotFound
	}

	return nil
}
