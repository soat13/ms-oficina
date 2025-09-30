package application

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/soat13/fase-1-oficina/internal/auth/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginService struct {
	repo   UserRepository
	tokens TokenService
}

func NewLoginService(repo UserRepository, tokens TokenService) *LoginService {
	return &LoginService{repo: repo, tokens: tokens}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

func (s *LoginService) Execute(ctx context.Context, in LoginInput) (LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	usernew, _ := domain.NewUser(user.ID, user.Name, user.Email, in.Password, user.CreatedAt)
	log.Println("pwd hash:", usernew.PasswordHash)

	if !user.VerifyPassword(in.Password) {
		return LoginOutput{}, ErrInvalidCredentials
	}
	tok, err := s.tokens.GenerateToken(user.ID, user.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	return LoginOutput{Token: tok}, nil
}
