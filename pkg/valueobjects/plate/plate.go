package plate

import (
	"regexp"
	"strings"
)

var (
	oldFormatRegex = regexp.MustCompile(`^[A-Z]{3}[0-9]{4}$`)
	newFormatRegex = regexp.MustCompile(`^[A-Z]{3}[0-9][A-Z][0-9]{2}$`)
)

type Plate struct {
	value string
}

func New(value string) (Plate, error) {
	cleanValue := strings.ToUpper(strings.TrimSpace(value))
	cleanValue = strings.ReplaceAll(cleanValue, " ", "")
	cleanValue = strings.ReplaceAll(cleanValue, "-", "")

	plate := Plate{value: cleanValue}

	if !plate.isValid() {
		return Plate{}, ErrInvalidPlate
	}

	return plate, nil
}

func (p Plate) String() string {
	return strings.ToUpper(p.value)
}

func (p Plate) isValid() bool {
	if p.value == "" {
		return false
	}
	return oldFormatRegex.MatchString(p.value) || newFormatRegex.MatchString(p.value)
}
