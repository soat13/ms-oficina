package application

import (
	"context"

	authDomain "github.com/soat13/fase-1-oficina/internal/auth/domain"
)

type (
	AuthenticateInput struct {
		CPF      string
		Password string
	}

	AuthenticateOutput struct {
		Token authDomain.Token
		User  UserView
	}

	AuthenticateUser struct {
		userReader   UserReader
		tokenService TokenService
	}
)

func NewAuthenticateUser(userReader UserReader, tokenService TokenService) *AuthenticateUser {
	return &AuthenticateUser{
		userReader:   userReader,
		tokenService: tokenService,
	}
}

func (uc *AuthenticateUser) Execute(ctx context.Context, in AuthenticateInput) (AuthenticateOutput, error) {
	user, err := uc.userReader.GetByCPF(ctx, in.CPF)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	if user == nil || !user.Password.Matches(in.Password) {
		return AuthenticateOutput{}, ErrInvalidCredentials
	}

	token, err := uc.tokenService.Generate(ctx, user)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	return AuthenticateOutput{Token: token, User: *user}, nil
}
