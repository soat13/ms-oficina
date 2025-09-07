package money

import "errors"

var (
	ErrMoneyNegative = errors.New("money cannot be negative")
)
