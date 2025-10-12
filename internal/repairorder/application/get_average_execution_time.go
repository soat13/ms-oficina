package application

import (
	"context"
)

type (
	GetAverageExecutionTimeOutput struct {
		AverageMinutes *float64
	}

	GetAverageExecutionTime struct {
		Repository Repository
	}
)

func NewGetAverageExecutionTime(repository Repository) *GetAverageExecutionTime {
	return &GetAverageExecutionTime{Repository: repository}
}

func (uc *GetAverageExecutionTime) Execute(ctx context.Context) (*GetAverageExecutionTimeOutput, error) {
	average, err := uc.Repository.GetAverageExecutionTime(ctx)
	if err != nil {
		return nil, err
	}

	return &GetAverageExecutionTimeOutput{AverageMinutes: average}, nil
}
