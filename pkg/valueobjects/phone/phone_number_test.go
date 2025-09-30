package phone

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    PhoneNumber
		wantErr error
	}{
		{
			name:    "valid mobile phone - 11 digits",
			value:   "81989017775",
			want:    PhoneNumber("81989017775"),
			wantErr: nil,
		},
		{
			name:    "valid mobile phone with formatting",
			value:   "(81) 98901-7775",
			want:    PhoneNumber("81989017775"),
			wantErr: nil,
		},
		{
			name:    "valid mobile phone with spaces",
			value:   "81 98901 7775",
			want:    PhoneNumber("81989017775"),
			wantErr: nil,
		},
		{
			name:    "valid phone with whitespace",
			value:   "  81989017775  ",
			want:    PhoneNumber("81989017775"),
			wantErr: nil,
		},
		{
			name:    "empty phone",
			value:   "",
			want:    PhoneNumber(""),
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name:    "phone with too many digits",
			value:   "819890177756",
			want:    PhoneNumber(""),
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name:    "phone with too few digits",
			value:   "8198901777",
			want:    PhoneNumber(""),
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name:    "phone with letters",
			value:   "8198901777a",
			want:    PhoneNumber(""),
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name:    "landline with 10 digits",
			value:   "1133334444",
			want:    PhoneNumber(""),
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name:    "phone with special characters",
			value:   "(81) 98901-7775",
			want:    PhoneNumber("81989017775"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.value)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, PhoneNumber(""), got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPhoneNumber_String(t *testing.T) {
	t.Parallel()

	phone, err := New("81989017775")
	require.NoError(t, err)

	got := phone.String()
	want := "81989017775"

	require.Equal(t, want, got)
}
