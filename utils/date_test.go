package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    time.Time
	}{
		{
			name:        "valid date YYYY-MM-DD",
			input:       "2020-05-15",
			expectError: false,
			expected:    time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "invalid format MM/DD/YYYY",
			input:       "05/15/2020",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid format",
			input:       "not-a-date",
			expectError: true,
		},
		{
			name:        "invalid date values",
			input:       "2020-13-45",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !result.Equal(tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}