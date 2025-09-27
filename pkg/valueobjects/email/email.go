package email

import "regexp"

type Email string

func New(value string) (Email, error) {
	email := Email(value)
	if !email.IsValid() {
		return Email(""), ErrInvalidEmail
	}
	return email, nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(string(e))
}
