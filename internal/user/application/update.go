package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/user/domain/role"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type (
	UpdateInput struct {
		ID          uuid.UUID
		Name        *string
		PhoneNumber *phone.PhoneNumber
		Email       *email.Email
		Password    *password.Password
		Roles       *role.Roles
	}

	UpdateUser struct {
		repo UserRepository
	}
)

func NewUpdateUser(repo UserRepository) *UpdateUser {
	return &UpdateUser{repo: repo}
}

func (uc *UpdateUser) Execute(ctx context.Context, in UpdateInput) error {
	user, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if in.Name != nil {
		if err := user.ChangeName(*in.Name); err != nil {
			return err
		}
	}

	if in.PhoneNumber != nil {
		user.ChangePhoneNumber(*in.PhoneNumber)
	}

	if in.Email != nil && user.Email.String() != in.Email.String() {
		exists, err := uc.repo.ExistsByEmail(ctx, in.Email.String())
		if err != nil {
			return err
		}
		if exists {
			return ErrDuplicateEmail
		}
		user.ChangeEmail(*in.Email)
	}

	if in.Password != nil {
		user.ChangePassword(*in.Password)
	}

	if in.Roles != nil {
		user.ChangeRoles(*in.Roles)
	}

	return uc.repo.Update(ctx, user)
}
