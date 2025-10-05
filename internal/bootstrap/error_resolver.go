package bootstrap

import (
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
)

func NewErrorResolver() *errorHelper.Resolver {
	errorResolver := errorHelper.NewErrorResolver()

	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidID)
	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidJSON)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)
	errorResolver.RegisterHTTPNotFoundError(sharedRepairOrder.ErrRepairOrderNotFound)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidDocument)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidPhoneNumber)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidEmail)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrPasswordTooShort)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrPasswordTooLong)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidPasswordHash)

	return errorResolver
}
