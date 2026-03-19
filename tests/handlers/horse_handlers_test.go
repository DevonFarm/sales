package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/DevonFarm/sales/horse"
)

// MockDB implements a simple in-memory database for testing
type MockDB struct {
	horses map[uuid.UUID]*horse.Horse
	farms  map[uuid.UUID]string // farmID -> farmName
	stats  map[uuid.UUID]*horse.DashboardStats
}

func NewMockDB() *MockDB {
	return &MockDB{
		horses: make(map[uuid.UUID]*horse.Horse),
		farms:  make(map[uuid.UUID]string),
		stats:  make(map[uuid.UUID]*horse.DashboardStats),
	}
}

// Mock QueryRow for horse creation
func (db *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) MockRow {
	if strings.Contains(sql, "INSERT INTO horses") {
		// Extract horse data from args
		name := args[0].(string)
		description := args[1].(string) 
		dateOfBirth := args[2].(time.Time)
		gender := args[3].(horse.Gender)
		farmID := args[4].(uuid.UUID)

		// Create new horse with generated ID
		horseID := uuid.New()
		h := &horse.Horse{
			ID:          horseID,
			Name:        name,
			Description: description,
			DateOfBirth: dateOfBirth,
			Gender:      gender,
			FarmID:      farmID,
		}
		db.horses[horseID] = h
		
		return MockRow{id: horseID}
	}
	
	return MockRow{err: fiber.ErrNotFound}
}

// Mock Query for horse retrieval
func (db *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (MockRows, error) {
	if strings.Contains(sql, "SELECT id, name, description, date_of_birth, gender FROM horses WHERE farm_id") {
		farmID := args[0].(uuid.UUID)
		var horses []*horse.Horse
		
		for _, h := range db.horses {
			if h.FarmID == farmID {
				horses = append(horses, h)
			}
		}
		
		return MockRows{horses: horses}, nil
	}
	
	if strings.Contains(sql, "SELECT COUNT(*) FROM horses WHERE farm_id") {
		farmID := args[0].(uuid.UUID)
		count := 0
		for _, h := range db.horses {
			if h.FarmID == farmID {
				count++
			}
		}
		return MockRows{count: count}, nil
	}
	
	if strings.Contains(sql, "SELECT gender, COUNT(*) FROM horses WHERE farm_id") {
		farmID := args[0].(uuid.UUID)
		genderCounts := make(map[horse.Gender]int)
		
		for _, h := range db.horses {
			if h.FarmID == farmID {
				genderCounts[h.Gender]++
			}
		}
		
		return MockRows{genderCounts: genderCounts}, nil
	}
	
	return MockRows{}, nil
}

type MockRow struct {
	id  uuid.UUID
	err error
}

func (r MockRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) > 0 {
		if id, ok := dest[0].(*uuid.UUID); ok {
			*id = r.id
		}
	}
	return nil
}

type MockRows struct {
	horses       []*horse.Horse
	count        int
	genderCounts map[horse.Gender]int
	index        int
	closed       bool
}

func (r *MockRows) Next() bool {
	if r.horses != nil {
		return r.index < len(r.horses)
	}
	if r.genderCounts != nil {
		return r.index < len(r.genderCounts)
	}
	return r.index == 0 && r.count > 0
}

func (r *MockRows) Scan(dest ...interface{}) error {
	if r.horses != nil && r.index < len(r.horses) {
		h := r.horses[r.index]
		if len(dest) >= 5 {
			*dest[0].(*uuid.UUID) = h.ID
			*dest[1].(*string) = h.Name
			*dest[2].(*string) = h.Description
			*dest[3].(*time.Time) = h.DateOfBirth
			*dest[4].(*horse.Gender) = h.Gender
		}
		r.index++
		return nil
	}
	
	if r.genderCounts != nil {
		genders := []horse.Gender{horse.GenderStallion, horse.GenderMare, horse.GenderGelding}
		if r.index < len(genders) {
			gender := genders[r.index]
			count := r.genderCounts[gender]
			*dest[0].(*horse.Gender) = gender
			*dest[1].(*int) = count
			r.index++
			return nil
		}
	}
	
	if r.count > 0 && r.index == 0 {
		*dest[0].(*int) = r.count
		r.index++
		return nil
	}
	
	return nil
}

func (r *MockRows) Close() {
	r.closed = true
}

func (r *MockRows) Err() error {
	return nil
}

func TestCreateHorse_ValidData(t *testing.T) {
	// Create mock database
	_ = NewMockDB()
	
	// Create Fiber app
	app := fiber.New()
	
	// Register the route (simplified version of the actual handler)
	app.Post("/farm/:farmID/horse", func(c *fiber.Ctx) error {
		var h horse.Horse
		if err := c.BodyParser(&h); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		
		// Parse date if provided
		dateStr := c.FormValue("date_of_birth")
		if dateStr != "" {
			dob, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			h.DateOfBirth = dob
		}
		
		// Parse farm ID
		farmIDStr := c.Params("farmID")
		farmID, err := uuid.Parse(farmIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid farm ID"})
		}
		h.FarmID = farmID
		
		// Validate gender
		if h.Gender.IsInvalid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid gender"})
		}
		
		// For this test, we'll simulate a successful save
		h.ID = uuid.New()
		
		return c.Status(fiber.StatusCreated).JSON(h)
	})
	
	// Test data
	farmID := uuid.New()
	horseData := map[string]interface{}{
		"name":         "Thunder",
		"description":  "A beautiful stallion",
		"gender":       horse.GenderStallion,
	}
	
	jsonData, _ := json.Marshal(horseData)
	
	// Create request with form data for date
	req := httptest.NewRequest("POST", "/farm/"+farmID.String()+"/horse?date_of_birth=2020-05-15", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	
	// Check response
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	
	// Parse response body
	var responseHorse horse.Horse
	err = json.NewDecoder(resp.Body).Decode(&responseHorse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Verify horse data
	if responseHorse.Name != "Thunder" {
		t.Errorf("Expected name 'Thunder', got '%s'", responseHorse.Name)
	}
	if responseHorse.Description != "A beautiful stallion" {
		t.Errorf("Expected description 'A beautiful stallion', got '%s'", responseHorse.Description)
	}
	if responseHorse.Gender != horse.GenderStallion {
		t.Errorf("Expected gender %v, got %v", horse.GenderStallion, responseHorse.Gender)
	}
	if responseHorse.FarmID != farmID {
		t.Errorf("Expected farm ID %s, got %s", farmID, responseHorse.FarmID)
	}
	if responseHorse.ID == uuid.Nil {
		t.Error("Expected horse ID to be set")
	}
}

func TestCreateHorse_InvalidGender(t *testing.T) {
	_ = NewMockDB()
	
	app := fiber.New()
	
	app.Post("/farm/:farmID/horse", func(c *fiber.Ctx) error {
		var h horse.Horse
		if err := c.BodyParser(&h); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		
		farmIDStr := c.Params("farmID")
		farmID, err := uuid.Parse(farmIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid farm ID"})
		}
		h.FarmID = farmID
		
		// Validate gender
		if h.Gender.IsInvalid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid gender"})
		}
		
		return c.Status(fiber.StatusCreated).JSON(h)
	})
	
	// Test data with invalid gender
	farmID := uuid.New()
	horseData := map[string]interface{}{
		"name":        "Invalid Horse",
		"description": "Horse with invalid gender",
		"gender":      99, // Invalid gender
	}
	
	jsonData, _ := json.Marshal(horseData)
	
	req := httptest.NewRequest("POST", "/farm/"+farmID.String()+"/horse", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	
	// Should return bad request
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
	
	// Parse error response
	var errorResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}
	
	if errorResp["error"] != "invalid gender" {
		t.Errorf("Expected error 'invalid gender', got '%s'", errorResp["error"])
	}
}

func TestCreateHorse_InvalidFarmID(t *testing.T) {
	app := fiber.New()
	
	app.Post("/farm/:farmID/horse", func(c *fiber.Ctx) error {
		var h horse.Horse
		if err := c.BodyParser(&h); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		
		farmIDStr := c.Params("farmID")
		_, err := uuid.Parse(farmIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid farm ID"})
		}
		
		return c.Status(fiber.StatusCreated).JSON(h)
	})
	
	// Test data
	horseData := map[string]interface{}{
		"name":        "Test Horse",
		"description": "Test description",
		"gender":      horse.GenderStallion,
	}
	
	jsonData, _ := json.Marshal(horseData)
	
	// Use invalid farm ID
	req := httptest.NewRequest("POST", "/farm/invalid-uuid/horse", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	
	// Should return bad request
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateHorse_MalformedJSON(t *testing.T) {
	app := fiber.New()
	
	app.Post("/farm/:farmID/horse", func(c *fiber.Ctx) error {
		var h horse.Horse
		if err := c.BodyParser(&h); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(h)
	})
	
	farmID := uuid.New()
	
	// Malformed JSON
	req := httptest.NewRequest("POST", "/farm/"+farmID.String()+"/horse", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	
	// Should return bad request
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}