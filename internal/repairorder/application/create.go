package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type CreateInput struct {
	CustomerID uuid.UUID
	VehicleID  uuid.UUID
}

type Create struct {
	repository     Repository
	customerReader CustomerReader
	vehicleReader  VehicleReader
}

func NewCreate(repository Repository, customerReader CustomerReader, vehicleReader VehicleReader) *Create {
	return &Create{
		repository:     repository,
		customerReader: customerReader,
		vehicleReader:  vehicleReader,
	}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) error {
	if err := uc.inputValidate(ctx, input); err != nil {
		return err
	}

	repairOrder, err := domain.NewRepairOrder(uuid.Nil, input.CustomerID, input.VehicleID, nil, nil, nil)
	if err != nil {
		return err
	}

	return uc.repository.Create(ctx, repairOrder)
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
