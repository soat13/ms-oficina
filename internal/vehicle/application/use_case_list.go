package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/maps"
	"github.com/soat13/oficina-utils/pkg/pagination"
)

type (
	ListInput struct {
		Pager pagination.Pagination
	}

	ListOutput struct {
		Vehicles []VehicleView
	}

	ListVehicles struct {
		repo VehicleRepository
	}
)

func NewListVehicles(repo VehicleRepository) *ListVehicles {
	return &ListVehicles{repo: repo}
}

func (lv *ListVehicles) Execute(ctx context.Context, in ListInput) (ListOutput, error) {
	vehicles, err := lv.repo.List(ctx, in.Pager)
	if err != nil {
		return ListOutput{}, err
	}

	return ListOutput{
		Vehicles: maps.Map(vehicles, toView),
	}, nil
}

type (
	ListByCustomerInput struct {
		CustomerID uuid.UUID
		Pager      pagination.Pagination
	}

	ListByCustomerOutput struct {
		Vehicles []VehicleView
	}

	ListVehiclesByCustomer struct {
		repo VehicleRepository
	}
)

func NewListVehiclesByCustomer(repo VehicleRepository) *ListVehiclesByCustomer {
	return &ListVehiclesByCustomer{repo: repo}
}

func (lvc *ListVehiclesByCustomer) Execute(ctx context.Context, in ListByCustomerInput) (ListByCustomerOutput, error) {
	vehicles, err := lvc.repo.ListByCustomerID(ctx, in.CustomerID, in.Pager)
	if err != nil {
		return ListByCustomerOutput{}, err
	}

	return ListByCustomerOutput{
		Vehicles: maps.Map(vehicles, toView),
	}, nil
}
