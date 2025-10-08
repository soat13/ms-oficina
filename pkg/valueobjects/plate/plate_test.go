package plate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid old format with dash",
			input:   "ABC-1234",
			want:    "ABC1234",
			wantErr: false,
		},
		{
			name:    "valid old format without dash",
			input:   "ABC1234",
			want:    "ABC1234",
			wantErr: false,
		},
		{
			name:    "valid new format",
			input:   "ABC1D23",
			want:    "ABC1D23",
			wantErr: false,
		},
		{
			name:    "valid with lowercase",
			input:   "abc1234",
			want:    "ABC1234",
			wantErr: false,
		},
		{
			name:    "valid with spaces",
			input:   " ABC 1234 ",
			want:    "ABC1234",
			wantErr: false,
		},
		{
			name:    "invalid format - too short",
			input:   "AB1234",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - too long",
			input:   "ABCD1234",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong pattern",
			input:   "123ABCD",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - mixed old and new",
			input:   "ABC1D234",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "only spaces",
			input:   "   ",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, ErrInvalidPlate, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}
