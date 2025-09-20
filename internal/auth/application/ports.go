package application

import (
    "context"

    "github.com/google/uuid"
    "github.com/soat13/fase-1-oficina/internal/auth/domain"
)

type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (domain.User, error)
    Create(ctx context.Context, u domain.User) error
}

type TokenService interface {
    GenerateToken(userID uuid.UUID, email string) (string, error)
    ParseToken(token string) (Claims, error)
}

type Claims struct {
    UserID uuid.UUID
    Email  string
}

