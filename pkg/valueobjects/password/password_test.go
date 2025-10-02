package password

import (
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
		wantErr       error
	}{
		{
			name:          "valid password",
			plainPassword: "mypassword123",
			wantErr:       nil,
		},
		{
			name:          "minimum valid password",
			plainPassword: "12345678",
			wantErr:       nil,
		},
		{
			name:          "password too short",
			plainPassword: "12345",
			wantErr:       ErrPasswordTooShort,
		},
		{
			name:          "password too long",
			plainPassword: strings.Repeat("a", 73),
			wantErr:       ErrPasswordTooLong,
		},
		{
			name:          "empty password",
			plainPassword: "",
			wantErr:       ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.plainPassword)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.Hash == "" {
				t.Error("New() returned empty hash for valid password")
			}
		})
	}
}

func TestPassword_Compare(t *testing.T) {
	plainPassword := "mypassword123"
	pwd, err := New(plainPassword)
	if err != nil {
		t.Fatalf("failed to create password: %v", err)
	}

	tests := []struct {
		name          string
		plainPassword string
		want          bool
	}{
		{
			name:          "correct password",
			plainPassword: "mypassword123",
			want:          true,
		},
		{
			name:          "incorrect password",
			plainPassword: "wrongpassword",
			want:          false,
		},
		{
			name:          "empty password",
			plainPassword: "",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pwd.Matches(tt.plainPassword)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromHash(t *testing.T) {
	validPassword, _ := New("testpassword123")
	validHash := validPassword.Hash

	tests := []struct {
		name    string
		hash    string
		wantErr error
	}{
		{
			name:    "valid hash",
			hash:    validHash,
			wantErr: nil,
		},
		{
			name:    "empty hash",
			hash:    "",
			wantErr: ErrInvalidPasswordHash,
		},
		{
			name:    "invalid hash",
			hash:    "not-a-valid-bcrypt-hash",
			wantErr: ErrInvalidPasswordHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromHash(tt.hash)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FromHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.Hash != tt.hash {
				t.Errorf("FromHash() hash = %v, want %v", got.Hash, tt.hash)
			}
		})
	}
}
