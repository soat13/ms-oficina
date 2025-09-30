package document

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    Document
		wantErr error
	}{
		{
			name:  "valid CPF",
			value: "11144477735",
			want: Document{
				Value:        "11144477735",
				DocumentType: CPF,
			},
			wantErr: nil,
		},
		{
			name:  "valid CNPJ",
			value: "11222333000181",
			want: Document{
				Value:        "11222333000181",
				DocumentType: CNPJ,
			},
			wantErr: nil,
		},
		{
			name:  "valid CPF with formatting",
			value: "111.444.777-35",
			want: Document{
				Value:        "11144477735",
				DocumentType: CPF,
			},
			wantErr: nil,
		},
		{
			name:  "valid CNPJ with formatting",
			value: "11.222.333/0001-81",
			want: Document{
				Value:        "11222333000181",
				DocumentType: CNPJ,
			},
			wantErr: nil,
		},
		{
			name:    "invalid document",
			value:   "12345678901",
			want:    Document{},
			wantErr: ErrInvalidDocument,
		},
		{
			name:    "empty string",
			value:   "",
			want:    Document{},
			wantErr: ErrInvalidDocument,
		},
		{
			name:    "whitespace only",
			value:   "   ",
			want:    Document{},
			wantErr: ErrInvalidDocument,
		},
		{
			name:  "CPF with whitespace",
			value: "  11144477735  ",
			want: Document{
				Value:        "11144477735",
				DocumentType: CPF,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.value)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, Document{}, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDocumentTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dt   DocumentType
		want string
	}{
		{
			name: "CPF type",
			dt:   CPF,
			want: "CPF",
		},
		{
			name: "CNPJ type",
			dt:   CNPJ,
			want: "CNPJ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, string(tt.dt))
		})
	}
}

func TestDocumentIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		doc  Document
		want bool
	}{
		{
			name: "valid CPF",
			doc:  Document{Value: "11144477735", DocumentType: CPF},
			want: true,
		},
		{
			name: "valid CNPJ",
			doc:  Document{Value: "11222333000181", DocumentType: CNPJ},
			want: true,
		},
		{
			name: "invalid document with type",
			doc:  Document{Value: "12345678901", DocumentType: CPF},
			want: false,
		},
		{
			name: "empty document",
			doc:  Document{},
			want: false,
		},
		{
			name: "document without type",
			doc:  Document{Value: "11144477735", DocumentType: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.doc.IsValid())
		})
	}
}
