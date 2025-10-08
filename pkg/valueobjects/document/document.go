package document

import (
	"strings"

	"github.com/paemuri/brdoc"
	string_helper "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
)

type DocumentType string

const (
	CPF  DocumentType = "CPF"
	CNPJ DocumentType = "CNPJ"
)

type Document struct {
	Value        string
	DocumentType DocumentType
}

func New(v string) (Document, error) {
	value := string_helper.OnlyNumbers(strings.TrimSpace(v))
	docType := getDocumentType(value)
	document := Document{
		Value:        value,
		DocumentType: docType,
	}

	if !document.IsValid() {
		return Document{}, ErrInvalidDocument
	}
	return document, nil
}

func (d Document) Type() string {
	return string(d.DocumentType)
}

func (d Document) IsValid() bool {
	return d.DocumentType != "" && (isCPF(d.Value) || isCNPJ(d.Value))
}

func getDocumentType(document string) DocumentType {
	document = strings.TrimSpace(document)
	switch {
	case isCPF(document):
		return CPF
	case isCNPJ(document):
		return CNPJ
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
