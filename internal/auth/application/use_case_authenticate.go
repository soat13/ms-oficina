package application

import (
	"context"

	authDomain "github.com/soat13/fase-1-oficina/internal/auth/domain"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
)

type (
	AuthenticateInput struct {
		Email    email.Email
		Password string
	}

	AuthenticateOutput struct {
		Token authDomain.Token
		User  UserView
	}

	AuthenticateUser struct {
		users  UserReader
		tokens TokenService
	}
)

func NewAuthenticateUser(users UserReader, tokens TokenService) *AuthenticateUser {
	return &AuthenticateUser{
		users:  users,
		tokens: tokens,
	}
}

func (uc *AuthenticateUser) Execute(ctx context.Context, in AuthenticateInput) (AuthenticateOutput, error) {
	user, err := uc.users.GetByEmail(ctx, in.Email.String())
	if err != nil {
		return AuthenticateOutput{}, err
	}

	if user == nil || !user.Password.Matches(in.Password) {
		return AuthenticateOutput{}, ErrInvalidCredentials
	}

	token, err := uc.tokens.Generate(ctx, user)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	return AuthenticateOutput{Token: token, User: *user}, nil
}
