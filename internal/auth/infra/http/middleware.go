package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	authApp "github.com/soat13/fase-1-oficina/internal/auth/application"
	authDomain "github.com/soat13/fase-1-oficina/internal/auth/domain"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
)

const (
	ClaimsContextKey = "auth.claims"
	UserIDContextKey = "auth.user_id"
)

type Middleware struct {
	tokens            authApp.TokenService
	errorHandler      *fiberHelper.ErrorHandler
	protectedPrefixes []string
}

func NewMiddleware(tokens authApp.TokenService, errorHandler *fiberHelper.ErrorHandler, prefixes []string) *Middleware {
	if len(prefixes) == 0 {
		prefixes = []string{"/admin"}
	}

	m := &Middleware{
		tokens:            tokens,
		errorHandler:      errorHandler,
		protectedPrefixes: prefixes,
	}

	m.errorHandler.ErrorResolver.RegisterHTTPUnauthorizedError(authDomain.ErrMissingToken)
	m.errorHandler.ErrorResolver.RegisterHTTPUnauthorizedError(authDomain.ErrInvalidToken)
	m.errorHandler.ErrorResolver.RegisterHTTPUnauthorizedError(authDomain.ErrExpiredToken)

	return m
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {
	if !m.shouldProtect(ctx.Path()) {
		return ctx.Next()
	}

	rawToken := ctx.Get("Authorization")
	if strings.TrimSpace(rawToken) == "" {
		return m.errorHandler.Handle(ctx, authDomain.ErrMissingToken)
	}

	token, err := extractBearerToken(rawToken)
	if err != nil {
		return m.errorHandler.Handle(ctx, authDomain.ErrInvalidToken)
	}

	claims, err := m.tokens.Validate(ctx.Context(), token)
	if err != nil {
		return m.errorHandler.Handle(ctx, err)
	}

	ctx.Locals(ClaimsContextKey, claims)
	ctx.Locals(UserIDContextKey, claims.UserID)

	return ctx.Next()
}

func (m *Middleware) shouldProtect(path string) bool {
	for _, prefix := range m.protectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func extractBearerToken(header string) (string, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return "", authDomain.ErrInvalidToken
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", authDomain.ErrInvalidToken
	}
	if strings.TrimSpace(parts[1]) == "" {
		return "", authDomain.ErrInvalidToken
	}
	return parts[1], nil
}
