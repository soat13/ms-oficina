package domain

import "errors"

var (
	ErrInvalidCustomerName         = errors.New("invalid customer name")
	ErrInvalidCustomerCellphone    = errors.New("invalid customer cellphone")
	ErrInvalidCustomerDocumentType = errors.New("invalid customer document type")
)
