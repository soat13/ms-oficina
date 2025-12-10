package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type (
	GetInput struct {
		ID uuid.UUID
	}

	GetOutput struct {
		RepairOrder RepairOrderView
	}

	GetRepairOrder struct {
		repo Repository
	}
)

func NewGetRepairOrder(repo Repository) *GetRepairOrder {
	return &GetRepairOrder{repo: repo}
}

func (uc *GetRepairOrder) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	repairOrder, err := uc.repo.GetById(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if repairOrder == nil {
		return nil, sharedRepairOrder.ErrRepairOrderNotFound
	}

	return &GetOutput{
		RepairOrder: toView(repairOrder),
	}, nil
}
