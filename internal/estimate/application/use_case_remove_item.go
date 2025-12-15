package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	RemoveItemInput struct {
		EstimateID uuid.UUID
		ItemID     uuid.UUID
	}

	RemoveItem struct {
		repository Repository
	}
)

func NewRemoveItem(repository Repository) *RemoveItem {
	return &RemoveItem{repository: repository}
}

func (r *RemoveItem) Execute(ctx context.Context, input RemoveItemInput) error {
	estimate, err := r.repository.GetByID(ctx, input.EstimateID)
	if err != nil {
		return err
	}
	if estimate == nil {
		return ErrEstimateNotFound
	}

	if err := estimate.RemoveItem(input.ItemID); err != nil {
		return err
	}

	return r.repository.Save(ctx, estimate)
}
