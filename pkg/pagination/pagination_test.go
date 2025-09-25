package pagination

import "testing"

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
			gotLimit, gotOffset := LimitAndOffset(tt.inputLimit, tt.inputOffset)

			if gotLimit != tt.wantLimit {
				t.Errorf("LimitAndOffset() limit = %v, want %v", gotLimit, tt.wantLimit)
			}

			if gotOffset != tt.wantOffset {
				t.Errorf("LimitAndOffset() offset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func TestLimitAndOffset_DefaultValue(t *testing.T) {
	limit, offset := LimitAndOffset(0, 0)

	expectedDefaultLimit := 50
	if limit != expectedDefaultLimit {
		t.Errorf("Default limit should be %d, got %d", expectedDefaultLimit, limit)
	}

	expectedDefaultOffset := 0
	if offset != expectedDefaultOffset {
		t.Errorf("Default offset should be %d, got %d", expectedDefaultOffset, offset)
	}
}
