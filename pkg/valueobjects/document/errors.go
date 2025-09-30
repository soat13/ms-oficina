package document

import "errors"

var (
	ErrInvalidDocument = errors.New("invalid document (must be valid CPF or CNPJ)")
)
