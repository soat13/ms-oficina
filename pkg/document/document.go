package document

import (
	"strings"

	"github.com/paemuri/brdoc"
)

func ValidateDocument(document string) error {
	document = strings.TrimSpace(document)
	if !isCPF(document) && !isCNPJ(document) {
		return ErrInvalidDocument
	}
	return nil
}

func GetDocumentType(document string) string {
	document = strings.TrimSpace(document)
	switch {
	case isCPF(document):
		return "CPF"
	case isCNPJ(document):
		return "CNPJ"
	default:
		return ""
	}
}

func isCPF(document string) bool {
	document = strings.TrimSpace(document)
	return brdoc.IsCPF(document)
}

func isCNPJ(document string) bool {
	document = strings.TrimSpace(document)
	return brdoc.IsCNPJ(document)
}
