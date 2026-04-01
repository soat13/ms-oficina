package application

import (
	"context"

	"github.com/google/uuid"
	authDomain "github.com/soat13/fase-1-oficina/internal/auth/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/authz"
	"github.com/soat13/oficina-utils/pkg/valueobjects/password"
)

type (
	UserView struct {
		ID       uuid.UUID
		Name     string
		Email    string
		Password password.Password
		Roles    authz.Roles
	}

	UserReader interface {
		GetByCPF(ctx context.Context, cpf string) (*UserView, error)
	}

	TokenService interface {
		Generate(ctx context.Context, user *UserView) (authDomain.Token, error)
		Validate(ctx context.Context, token string) (authDomain.Claims, error)
	}
)
