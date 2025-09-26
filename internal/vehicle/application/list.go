package application

import (
	"context"
)

type ListInput struct {
	Limit  int
	Offset int
}

type ListVehicles struct {
	repo VehicleRepository
}

func NewListVehicles(repo VehicleRepository) *ListVehicles {
	return &ListVehicles{repo: repo}
}

func (uc *ListVehicles) Execute(ctx context.Context, in ListInput) ([]VehicleView, error) {
	vehicles, err := uc.repo.List(ctx, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	var views []VehicleView
	for _, v := range vehicles {
		views = append(views, toView(v))
	}

	return views, nil
}
