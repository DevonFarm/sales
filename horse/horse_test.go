package horse

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHorse_Age(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name        string
		dateOfBirth time.Time
		expected    int
	}{
		{
			name:        "1 year old",
			dateOfBirth: now.AddDate(-1, 0, 0),
			expected:    1,
		},
		{
			name:        "5 years old",
			dateOfBirth: now.AddDate(-5, 0, 0),
			expected:    5,
		},
		{
			name:        "newborn",
			dateOfBirth: now.AddDate(0, 0, -1), // 1 day old
			expected:    0,
		},
		{
			name:        "10 years old",
			dateOfBirth: now.AddDate(-10, -6, 0), // 10.5 years, should round down
			expected:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			horse := &Horse{
				DateOfBirth: tt.dateOfBirth,
			}
			result := horse.Age()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestHorse_HTMLPath(t *testing.T) {
	tests := []struct {
		name     string
		horseName string
		expected string
	}{
		{
			name:      "simple name",
			horseName: "Thunder",
			expected:  "thunder.html",
		},
		{
			name:      "name with spaces",
			horseName: "Black Beauty",
			expected:  "black_beauty.html",
		},
		{
			name:      "name with special chars",
			horseName: "Mr. Ed's Horse",
			expected:  "mr__ed's_horse.html",
		},
		{
			name:      "camelCase name",
			horseName: "ThunderBolt",
			expected:  "thunder_bolt.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			horse := &Horse{
				Name: tt.horseName,
			}
			result := horse.HTMLPath()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHorse_NewImage(t *testing.T) {
	horse := &Horse{
		Name: "Thunder",
	}

	// Test adding first image
	horse.NewImage("full.jpg", "thumb.jpg", "Thunder standing")
	
	if len(horse.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(horse.Images))
	}

	img := horse.Images[0]
	expectedFull := "assets/images/horses/Thunder/full.jpg"
	expectedThumb := "assets/images/horses/Thunder/thumb.jpg"
	expectedAlt := "Thunder standing"

	if img.Full != expectedFull {
		t.Errorf("expected full path %s, got %s", expectedFull, img.Full)
	}
	if img.Thumbnail != expectedThumb {
		t.Errorf("expected thumbnail path %s, got %s", expectedThumb, img.Thumbnail)
	}
	if img.Alt != expectedAlt {
		t.Errorf("expected alt text %s, got %s", expectedAlt, img.Alt)
	}

	// Test adding second image
	horse.NewImage("side.jpg", "side_thumb.jpg", "Thunder from side")
	
	if len(horse.Images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(horse.Images))
	}
}

func TestHorse_ValidateGender(t *testing.T) {
	tests := []struct {
		name          string
		gender        Gender
		expectError   bool
	}{
		{
			name:        "valid stallion",
			gender:      GenderStallion,
			expectError: false,
		},
		{
			name:        "valid mare", 
			gender:      GenderMare,
			expectError: false,
		},
		{
			name:        "valid gelding",
			gender:      GenderGelding,
			expectError: false,
		},
		{
			name:        "invalid gender",
			gender:      Gender(99),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			horse := &Horse{
				ID:          uuid.New(),
				Name:        "Test Horse",
				Description: "Test Description",
				DateOfBirth: time.Now().AddDate(-5, 0, 0),
				Gender:      tt.gender,
				FarmID:      uuid.New(),
			}

			// We can't easily test Save() without a database connection,
			// but we can test the validation logic by checking IsInvalid()
			hasError := horse.Gender.IsInvalid()
			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}