package phone

import (
	"strings"

	stringHelper "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
)

type PhoneNumber string

func New(v string) (PhoneNumber, error) {
	value := stringHelper.OnlyNumbers(strings.TrimSpace(v))
	phone := PhoneNumber(value)

	if !phone.isValid() {
		return PhoneNumber(""), ErrInvalidPhoneNumber
	}

	return phone, nil
}

func (p PhoneNumber) String() string {
	return string(p)
}

func (p PhoneNumber) isValid() bool {
	return len(string(p)) == 11
}
