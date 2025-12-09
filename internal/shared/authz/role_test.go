package authz

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Role
		wantErr error
	}{
		{
			name:    "valid attendant",
			value:   "attendant",
			want:    Role(attendant),
			wantErr: nil,
		},
		{
			name:    "valid manager",
			value:   "manager",
			want:    Role(manager),
			wantErr: nil,
		},
		{
			name:    "valid mechanic",
			value:   "mechanic",
			want:    Role(mechanic),
			wantErr: nil,
		},
		{
			name:    "case insensitive",
			value:   "ATTENDANT",
			want:    Role(attendant),
			wantErr: nil,
		},
		{
			name:    "with whitespace",
			value:   "  manager  ",
			want:    Role(manager),
			wantErr: nil,
		},
		{
			name:    "invalid role",
			value:   "invalid_role",
			want:    Role(""),
			wantErr: ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.value)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got != tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRoles(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		want    Roles
		wantErr error
	}{
		{
			name:    "valid single role",
			values:  []string{"attendant"},
			want:    Roles{Role(attendant)},
			wantErr: nil,
		},
		{
			name:    "valid multiple roles",
			values:  []string{"attendant", "manager", "mechanic"},
			want:    Roles{Role(attendant), Role(manager), Role(mechanic)},
			wantErr: nil,
		},
		{
			name:    "case insensitive",
			values:  []string{"ATTENDANT", "Manager", "mEcHaNiC"},
			want:    Roles{Role(attendant), Role(manager), Role(mechanic)},
			wantErr: nil,
		},
		{
			name:    "with whitespace",
			values:  []string{"  attendant  ", " manager "},
			want:    Roles{Role(attendant), Role(manager)},
			wantErr: nil,
		},
		{
			name:    "empty list",
			values:  []string{},
			want:    nil,
			wantErr: ErrRolesRequired,
		},
		{
			name:    "nil list",
			values:  nil,
			want:    nil,
			wantErr: ErrRolesRequired,
		},
		{
			name:    "invalid role",
			values:  []string{"invalid_role"},
			want:    nil,
			wantErr: ErrInvalidRole,
		},
		{
			name:    "mixed valid and invalid",
			values:  []string{"attendant", "invalid_role"},
			want:    nil,
			wantErr: ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRoles(tt.values)
			if err != tt.wantErr {
				t.Errorf("NewRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if len(got) != len(tt.want) {
					t.Errorf("NewRoles() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i, role := range got {
					if role != tt.want[i] {
						t.Errorf("NewRoles()[%d] = %v, want %v", i, role, tt.want[i])
					}
				}
			}
		})
	}
}

func TestRolesStrings(t *testing.T) {
	tests := []struct {
		name  string
		roles Roles
		want  []string
	}{
		{
			name:  "single role",
			roles: Roles{Role(attendant)},
			want:  []string{"attendant"},
		},
		{
			name:  "multiple roles",
			roles: Roles{Role(attendant), Role(manager), Role(mechanic)},
			want:  []string{"attendant", "manager", "mechanic"},
		},
		{
			name:  "empty roles",
			roles: Roles{},
			want:  []string{},
		},
		{
			name:  "nil roles",
			roles: nil,
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.roles.Strings()
			if len(got) != len(tt.want) {
				t.Errorf("Roles.Strings() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i, str := range got {
				if str != tt.want[i] {
					t.Errorf("Roles.Strings()[%d] = %v, want %v", i, str, tt.want[i])
				}
			}
		})
	}
}
