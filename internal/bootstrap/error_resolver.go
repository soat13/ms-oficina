package bootstrap

import (
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/soat13/oficina-utils/pkg/db/bun_helper"
	errorHelper "github.com/soat13/oficina-utils/pkg/error"
	"github.com/soat13/oficina-utils/pkg/http/fiber"
	"github.com/soat13/oficina-utils/pkg/valueobjects/email"
)

func NewErrorResolver() *errorHelper.Resolver {
	errorResolver := errorHelper.NewErrorResolver()

	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidID)
	errorResolver.RegisterHTTPBadRequestError(fiber.ErrInvalidID)
	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidJSON)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)
	errorResolver.RegisterHTTPConflictError(bun_helper.ErrResourceInUse)
	errorResolver.RegisterHTTPNotFoundError(sharedRepairOrder.ErrRepairOrderNotFound)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidDocument)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidPhoneNumber)
	errorResolver.RegisterHTTPUnprocessableError(email.ErrInvalidEmail)

	return errorResolver
}
