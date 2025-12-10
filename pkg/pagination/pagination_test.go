package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitAndOffset(t *testing.T) {
	tests := []struct {
		name        string
		inputLimit  int
		inputOffset int
		wantLimit   int
		wantOffset  int
	}{
		{
			name:        "valid positive values",
			inputLimit:  10,
			inputOffset: 20,
			wantLimit:   10,
			wantOffset:  20,
		},
		{
			name:        "zero limit should use default",
			inputLimit:  0,
			inputOffset: 10,
			wantLimit:   50,
			wantOffset:  10,
		},
		{
			name:        "negative limit should use default",
			inputLimit:  -5,
			inputOffset: 10,
			wantLimit:   50,
			wantOffset:  10,
		},
		{
			name:        "zero offset should remain zero",
			inputLimit:  25,
			inputOffset: 0,
			wantLimit:   25,
			wantOffset:  0,
		},
		{
			name:        "negative offset should become zero",
			inputLimit:  25,
			inputOffset: -10,
			wantLimit:   25,
			wantOffset:  0,
		},
		{
			name:        "both invalid values",
			inputLimit:  -1,
			inputOffset: -5,
			wantLimit:   50,
			wantOffset:  0,
		},
		{
			name:        "large valid values",
			inputLimit:  1000,
			inputOffset: 5000,
			wantLimit:   1000,
			wantOffset:  5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pager := New(tt.inputLimit, tt.inputOffset)
			assert.Equal(t, pager.Limit, tt.wantLimit)
			assert.Equal(t, pager.Offset, tt.wantOffset)
		})
	}
}

func TestLimitAndOffset_DefaultValue(t *testing.T) {
	pager := New(0, 0)
	expectedDefaultLimit := 50
	expectedDefaultOffset := 0
	assert.Equal(t, pager.Limit, expectedDefaultLimit)
	assert.Equal(t, pager.Offset, expectedDefaultOffset)
}

func TestAtoiDefault(t *testing.T) {
	tests := []struct {
		name  string
		input string
		def   int
		want  int
	}{
		{name: "valid input", input: "10", def: 50, want: 10},
		{name: "empty input", input: "", def: 50, want: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AtoiDefault(tt.input, tt.def)
			if got != tt.want {
				t.Errorf("AtoiDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
