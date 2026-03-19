package horse

import (
	"testing"
	"time"
)

func TestGender_IsInvalid(t *testing.T) {
	tests := []struct {
		name     string
		gender   Gender
		expected bool
	}{
		{
			name:     "valid stallion",
			gender:   GenderStallion,
			expected: false,
		},
		{
			name:     "valid mare",
			gender:   GenderMare,
			expected: false,
		},
		{
			name:     "valid gelding",
			gender:   GenderGelding,
			expected: false,
		},
		{
			name:     "invalid negative",
			gender:   Gender(-1),
			expected: true,
		},
		{
			name:     "invalid too high",
			gender:   Gender(10),
			expected: true,
		},
		{
			name:     "invalid zero (GenderInvalid)",
			gender:   GenderInvalid,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gender.IsInvalid()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHorse_GenderString(t *testing.T) {
	tests := []struct {
		name     string
		horse    *Horse
		expected string
	}{
		{
			name: "adult stallion",
			horse: &Horse{
				Gender:      GenderStallion,
				DateOfBirth: time.Now().AddDate(-5, 0, 0), // 5 years old
			},
			expected: "Stallion",
		},
		{
			name: "young stallion (colt)",
			horse: &Horse{
				Gender:      GenderStallion,
				DateOfBirth: time.Now().AddDate(-2, 0, 0), // 2 years old
			},
			expected: "Colt",
		},
		{
			name: "adult mare",
			horse: &Horse{
				Gender:      GenderMare,
				DateOfBirth: time.Now().AddDate(-6, 0, 0), // 6 years old
			},
			expected: "Mare",
		},
		{
			name: "young mare (filly)",
			horse: &Horse{
				Gender:      GenderMare,
				DateOfBirth: time.Now().AddDate(-1, 0, 0), // 1 year old
			},
			expected: "Filly",
		},
		{
			name: "adult gelding",
			horse: &Horse{
				Gender:      GenderGelding,
				DateOfBirth: time.Now().AddDate(-8, 0, 0), // 8 years old
			},
			expected: "Gelding",
		},
		{
			name: "young gelding (still gelding, not colt)",
			horse: &Horse{
				Gender:      GenderGelding,
				DateOfBirth: time.Now().AddDate(-3, 0, 0), // 3 years old
			},
			expected: "Gelding",
		},
		{
			name: "invalid gender",
			horse: &Horse{
				Gender:      Gender(99),
				DateOfBirth: time.Now().AddDate(-5, 0, 0),
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.horse.GenderString()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}