package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    Email
		wantErr error
	}{
		{
			name:    "valid email",
			value:   "test@example.com",
			want:    Email("test@example.com"),
			wantErr: nil,
		},
		{
			name:    "valid email with subdomain",
			value:   "user@mail.example.com",
			want:    Email("user@mail.example.com"),
			wantErr: nil,
		},
		{
			name:    "valid email with numbers",
			value:   "user123@example123.com",
			want:    Email("user123@example123.com"),
			wantErr: nil,
		},
		{
			name:    "valid email with special chars",
			value:   "user.name+tag@example.com",
			want:    Email("user.name+tag@example.com"),
			wantErr: nil,
		},
		{
			name:    "empty email",
			value:   "",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email without @",
			value:   "testexample.com",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email without domain",
			value:   "test@",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email without local part",
			value:   "@example.com",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email with spaces",
			value:   "test @example.com",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email with invalid domain",
			value:   "test@example",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email with invalid characters",
			value:   "test@exam ple.com",
			want:    Email(""),
			wantErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.value)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, Email(""), got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	t.Parallel()

	email, err := New("test@example.com")
	require.NoError(t, err)

	got := email.String()
	want := "test@example.com"

	require.Equal(t, want, got)
}
