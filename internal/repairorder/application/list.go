package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type (
	ListInput struct {
		Pager pagination.Pagination
	}

	ListOutput struct {
		RepairOrders []RepairOrderView
	}

	ListRepairOrders struct {
		repo Repository
	}
)

func NewListRepairOrders(repo Repository) *ListRepairOrders {
	return &ListRepairOrders{repo: repo}
}

func (uc *ListRepairOrders) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	items, err := uc.repo.List(ctx, in.Pager)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		RepairOrders: maps.Map(items, toView),
	}, nil
}
