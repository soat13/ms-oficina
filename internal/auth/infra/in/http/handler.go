package http

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	authApp "github.com/soat13/fase-1-oficina/internal/auth/application"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
)

type Handler struct {
	authenticate *authApp.AuthenticateUser
	validator    *validator.Validate
	errorHandler *fiberHelper.ErrorHandler
}

func NewHandler(
	authenticate *authApp.AuthenticateUser,
	validator *validator.Validate,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	handler := &Handler{
		authenticate: authenticate,
		validator:    validator,
		errorHandler: errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPUnauthorizedError(authApp.ErrInvalidCredentials)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	app.Post("/auth/login", h.login)
}

type loginBody struct {
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type authenticatedUserJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Roles []string  `json:"roles"`
}

type loginResponse struct {
	AccessToken string                `json:"access_token"`
	TokenType   string                `json:"token_type"`
	ExpiresAt   time.Time             `json:"expires_at"`
	User        authenticatedUserJSON `json:"user"`
}

func (h *Handler) login(ctx *fiber.Ctx) error {
	var body loginBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	emailVO, err := email.New(body.Email)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	out, err := h.authenticate.Execute(ctx.Context(), authApp.AuthenticateInput{
		Email:    emailVO,
		Password: body.Password,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	resp := loginResponse{
		AccessToken: out.Token.Value,
		TokenType:   "Bearer",
		ExpiresAt:   out.Token.ExpiresAt,
		User: authenticatedUserJSON{
			ID:    out.User.ID,
			Name:  out.User.Name,
			Email: out.User.Email,
			Roles: out.User.Roles.Strings(),
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
