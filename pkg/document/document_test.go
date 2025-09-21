package document

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDocument(t *testing.T) {
	t.Parallel()

	t.Run("should return nil for valid CPF", func(t *testing.T) {
		err := ValidateDocument("11144477735")
		require.NoError(t, err)
	})

	t.Run("should return nil for valid CNPJ", func(t *testing.T) {
		err := ValidateDocument("11222333000181")
		require.NoError(t, err)
	})

	t.Run("should return error for invalid document", func(t *testing.T) {
		err := ValidateDocument("12345678901")
		require.ErrorIs(t, err, ErrInvalidDocument)
	})
}

func TestGetDocumentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		document string
		want     string
	}{
		{"valid CPF", "11144477735", "CPF"},
		{"valid CNPJ", "11222333000181", "CNPJ"},
		{"invalid document", "12345678901", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDocumentType(tt.document)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestIsCPF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		document string
		want     bool
	}{
		{"valid CPF", "11144477735", true},
		{"valid CPF with dots", "111.444.777-35", true},
		{"invalid CPF", "12345678901", false},
		{"empty string", "", false},
		{"whitespace", "   ", false},
		{"CPF with whitespace", "  11144477735  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCPF(tt.document)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestIsCNPJ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		document string
		want     bool
	}{
		{"valid CNPJ", "11222333000181", true},
		{"valid CNPJ with dots", "11.222.333/0001-81", true},
		{"invalid CNPJ", "12345678901234", false},
		{"empty string", "", false},
		{"whitespace", "   ", false},
		{"CNPJ with whitespace", "  11222333000181  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCNPJ(tt.document)
			require.Equal(t, tt.want, result)
		})
	}
}
