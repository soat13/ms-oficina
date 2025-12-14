package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
	"github.com/soat13/fase-1-oficina/pkg/pagination"
)

type (
	ProductView struct {
		ID    uuid.UUID
		Name  string
		Price money.Money
		Stock int
	}

	Repository interface {
		Create(ctx context.Context, product *domain.Product) error
		Update(ctx context.Context, product *domain.Product) error
		UpdateBatch(ctx context.Context, products []*domain.Product) error
		Delete(ctx context.Context, id uuid.UUID) error
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error)
		List(ctx context.Context, pager pagination.Pagination) ([]*domain.Product, error)
		ExistsByName(ctx context.Context, name string) (bool, error)
	}

	EventPublisher interface {
		PublishStockInsufficientDetected(ctx context.Context, estimateID uuid.UUID, repairOrderID uuid.UUID) error
		PublishStockReduceConfirmed(ctx context.Context, estimateID uuid.UUID, repairOrderID uuid.UUID) error
	}
)

func toView(p *domain.Product) ProductView {
	return ProductView{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}
}
