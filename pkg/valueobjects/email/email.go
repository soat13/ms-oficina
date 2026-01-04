package email

import (
	"regexp"
	"strings"
)

type Email string

func New(value string) (Email, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	email := Email(normalized)
	if !email.IsValid() {
		return Email(""), ErrInvalidEmail
	}
	return email, nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	emailRegex := `^[a-zA-Z0-9](?:[a-zA-Z0-9._%+-]{0,62}[a-zA-Z0-9])?@[A-Za-z0-9](?:[A-Za-z0-9-]{0,62}[A-Za-z0-9])?(?:\.[A-Za-z]{2,})+$`
	return regexp.MustCompile(emailRegex).MatchString(string(e))
}
