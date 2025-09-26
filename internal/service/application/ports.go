package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type ServiceView struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Price money.Money `json:"price"`
}

func toView(s *domain.Service) ServiceView {
	return ServiceView{
		ID:    s.ID,
		Name:  s.Name,
		Price: s.Price,
	}
}

type ServiceRepository interface {
	Create(ctx context.Context, s *domain.Service) error
	Update(ctx context.Context, s *domain.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Service, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
