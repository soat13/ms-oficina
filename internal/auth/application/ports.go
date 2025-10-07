package application

import (
	"context"

	authDomain "github.com/soat13/fase-1-oficina/internal/auth/domain"
	userDomain "github.com/soat13/fase-1-oficina/internal/user/domain"
)

type (
	UserReader interface {
		GetByEmail(ctx context.Context, email string) (*userDomain.User, error)
	}

	TokenService interface {
		Generate(ctx context.Context, user *userDomain.User) (authDomain.Token, error)
		Validate(ctx context.Context, token string) (authDomain.Claims, error)
	}
)
