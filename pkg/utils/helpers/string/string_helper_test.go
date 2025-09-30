package string_helper

import (
	"testing"
)

func TestOnlyNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CPF with formatting",
			input:    "123.456.789-01",
			expected: "12345678901",
		},
		{
			name:     "CNPJ with formatting",
			input:    "12.345.678/0001-90",
			expected: "12345678000190",
		},
		{
			name:     "only numbers",
			input:    "12345678901",
			expected: "12345678901",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only letters",
			input:    "abcdef",
			expected: "",
		},
		{
			name:     "mixed characters",
			input:    "abc123def456",
			expected: "123456",
		},
		{
			name:     "special characters",
			input:    "!@#123$%^456&*()",
			expected: "123456",
		},
		{
			name:     "spaces and tabs",
			input:    " 1 2 3 	4 5 6 ",
			expected: "123456",
		},
		{
			name:     "phone number with formatting",
			input:    "(11) 98765-4321",
			expected: "11987654321",
		},
		{
			name:     "numbers with letters and symbols",
			input:    "ABC123-DEF456.GHI789",
			expected: "123456789",
		},
		{
			name:     "unicode characters",
			input:    "123ção456",
			expected: "123456",
		},
		{
			name:     "leading and trailing non-numbers",
			input:    "abc123def",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OnlyNumbers(tt.input)
			if result != tt.expected {
				t.Errorf("OnlyNumbers(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
