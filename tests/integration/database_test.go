package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/farm"
	"github.com/DevonFarm/sales/horse"
	"github.com/DevonFarm/sales/user"
)

var testDB *database.DB

func TestMain(m *testing.M) {
	// Setup test database
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		fmt.Printf("Failed to setup test database: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if testDB != nil {
		testDB.Close(context.Background())
	}

	os.Exit(code)
}

func setupTestDB() (*database.DB, error) {
	// Use test database connection string
	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == "" {
		// If no test database URL, skip integration tests
		return nil, fmt.Errorf("TEST_DATABASE_URL not set - skipping integration tests")
	}

	db, err := database.NewDBConn(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Verify tables exist (migrations should be run separately)
	ctx := context.Background()
	tables := []string{"users", "farms", "horses"}
	for _, table := range tables {
		var exists bool
		err := db.QueryRow(ctx, 
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", 
			table,
		).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return nil, fmt.Errorf("table %s does not exist - run migrations first", table)
		}
	}

	return db, nil
}

func cleanupTestData(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	ctx := context.Background()
	
	// Clean up in reverse dependency order
	_, err := testDB.Exec(ctx, "DELETE FROM horses")
	if err != nil {
		t.Fatalf("Failed to cleanup horses: %v", err)
	}
	
	_, err = testDB.Exec(ctx, "DELETE FROM farms")
	if err != nil {
		t.Fatalf("Failed to cleanup farms: %v", err)
	}
	
	_, err = testDB.Exec(ctx, "DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}

func TestUserCRUD(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}
	
	defer cleanupTestData(t)
	
	ctx := context.Background()

	// Test user creation
	testUser, err := user.NewUser(ctx, testDB, "Test User", "test@example.com", "stytch-123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if testUser.ID == uuid.Nil {
		t.Error("User ID should not be nil after creation")
	}

	// Test user retrieval by ID
	retrievedUser, err := user.GetUser(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.Name != testUser.Name {
		t.Errorf("Expected name %s, got %s", testUser.Name, retrievedUser.Name)
	}
	if retrievedUser.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, retrievedUser.Email)
	}
	if retrievedUser.StytchID != testUser.StytchID {
		t.Errorf("Expected stytch ID %s, got %s", testUser.StytchID, retrievedUser.StytchID)
	}

	// Test user retrieval by Stytch ID
	userByStytch, err := user.GetUserByStytchID(ctx, testDB, testUser.StytchID)
	if err != nil {
		t.Fatalf("Failed to get user by Stytch ID: %v", err)
	}

	if userByStytch.ID != testUser.ID {
		t.Errorf("Expected user ID %s, got %s", testUser.ID, userByStytch.ID)
	}

	// Test user update
	testUser.Name = "Updated Name"
	testUser.Email = "updated@example.com"
	err = testUser.Update(ctx, testDB)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	updatedUser, err := user.GetUser(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updatedUser.Name != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got %s", updatedUser.Name)
	}
	if updatedUser.Email != "updated@example.com" {
		t.Errorf("Expected updated email 'updated@example.com', got %s", updatedUser.Email)
	}
}

func TestFarmCRUD(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}
	
	defer cleanupTestData(t)
	
	ctx := context.Background()

	// Create a test user first
	testUser, err := user.NewUser(ctx, testDB, "Farm Owner", "owner@example.com", "stytch-456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test farm creation
	testFarm := &farm.Farm{
		Name: "Test Farm",
	}

	err = testFarm.Save(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to create farm: %v", err)
	}

	if testFarm.ID == uuid.Nil {
		t.Error("Farm ID should not be nil after creation")
	}

	// Test farm retrieval
	retrievedFarm, err := farm.GetFarm(ctx, testDB, testFarm.ID.String())
	if err != nil {
		t.Fatalf("Failed to get farm: %v", err)
	}

	if retrievedFarm.Name != testFarm.Name {
		t.Errorf("Expected farm name %s, got %s", testFarm.Name, retrievedFarm.Name)
	}

	// Verify user was updated with farm ID
	updatedUser, err := user.GetUser(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to get user after farm creation: %v", err)
	}

	if updatedUser.FarmID != testFarm.ID {
		t.Errorf("Expected user farm ID %s, got %s", testFarm.ID, updatedUser.FarmID)
	}
}

func TestHorseCRUD(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}
	
	defer cleanupTestData(t)
	
	ctx := context.Background()

	// Create test user and farm
	testUser, err := user.NewUser(ctx, testDB, "Horse Owner", "horses@example.com", "stytch-789")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	testFarm := &farm.Farm{
		Name: "Horse Farm",
	}

	err = testFarm.Save(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to create test farm: %v", err)
	}

	// Test horse creation
	testHorse := &horse.Horse{
		Name:        "Thunder",
		Description: "A beautiful stallion",
		DateOfBirth: time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:      horse.GenderStallion,
		FarmID:      testFarm.ID,
	}

	err = testHorse.Save(ctx, testDB)
	if err != nil {
		t.Fatalf("Failed to create horse: %v", err)
	}

	if testHorse.ID == uuid.Nil {
		t.Error("Horse ID should not be nil after creation")
	}

	// Test horse retrieval by farm
	horses, err := horse.GetHorsesByFarmID(ctx, testDB, testFarm.ID)
	if err != nil {
		t.Fatalf("Failed to get horses by farm ID: %v", err)
	}

	if len(horses) != 1 {
		t.Fatalf("Expected 1 horse, got %d", len(horses))
	}

	retrievedHorse := horses[0]
	if retrievedHorse.Name != testHorse.Name {
		t.Errorf("Expected horse name %s, got %s", testHorse.Name, retrievedHorse.Name)
	}
	if retrievedHorse.Description != testHorse.Description {
		t.Errorf("Expected horse description %s, got %s", testHorse.Description, retrievedHorse.Description)
	}
	if retrievedHorse.Gender != testHorse.Gender {
		t.Errorf("Expected horse gender %v, got %v", testHorse.Gender, retrievedHorse.Gender)
	}

	// Test dashboard stats
	stats, err := horse.GetDashboardStats(ctx, testDB, testFarm.ID)
	if err != nil {
		t.Fatalf("Failed to get dashboard stats: %v", err)
	}

	if stats.TotalHorses != 1 {
		t.Errorf("Expected 1 total horse, got %d", stats.TotalHorses)
	}
	if stats.Stallions != 1 {
		t.Errorf("Expected 1 stallion, got %d", stats.Stallions)
	}
	if stats.Mares != 0 {
		t.Errorf("Expected 0 mares, got %d", stats.Mares)
	}
	if stats.Geldings != 0 {
		t.Errorf("Expected 0 geldings, got %d", stats.Geldings)
	}

	// Test adding multiple horses for stats
	mare := &horse.Horse{
		Name:        "Beauty",
		Description: "A beautiful mare",
		DateOfBirth: time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC),
		Gender:      horse.GenderMare,
		FarmID:      testFarm.ID,
	}

	err = mare.Save(ctx, testDB)
	if err != nil {
		t.Fatalf("Failed to create mare: %v", err)
	}

	gelding := &horse.Horse{
		Name:        "Steady",
		Description: "A calm gelding",
		DateOfBirth: time.Date(2017, 8, 22, 0, 0, 0, 0, time.UTC),
		Gender:      horse.GenderGelding,
		FarmID:      testFarm.ID,
	}

	err = gelding.Save(ctx, testDB)
	if err != nil {
		t.Fatalf("Failed to create gelding: %v", err)
	}

	// Test updated stats
	updatedStats, err := horse.GetDashboardStats(ctx, testDB, testFarm.ID)
	if err != nil {
		t.Fatalf("Failed to get updated dashboard stats: %v", err)
	}

	if updatedStats.TotalHorses != 3 {
		t.Errorf("Expected 3 total horses, got %d", updatedStats.TotalHorses)
	}
	if updatedStats.Stallions != 1 {
		t.Errorf("Expected 1 stallion, got %d", updatedStats.Stallions)
	}
	if updatedStats.Mares != 1 {
		t.Errorf("Expected 1 mare, got %d", updatedStats.Mares)
	}
	if updatedStats.Geldings != 1 {
		t.Errorf("Expected 1 gelding, got %d", updatedStats.Geldings)
	}
}

func TestDatabaseConstraints(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}
	
	defer cleanupTestData(t)
	
	ctx := context.Background()

	// Test invalid horse gender constraint
	testUser, err := user.NewUser(ctx, testDB, "Constraint Tester", "constraint@example.com", "stytch-constraint")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	testFarm := &farm.Farm{
		Name: "Constraint Farm",
	}

	err = testFarm.Save(ctx, testDB, testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to create test farm: %v", err)
	}

	// Try to create horse with invalid gender
	invalidHorse := &horse.Horse{
		Name:        "Invalid",
		Description: "Invalid gender horse",
		DateOfBirth: time.Now(),
		Gender:      horse.Gender(99), // Invalid gender
		FarmID:      testFarm.ID,
	}

	err = invalidHorse.Save(ctx, testDB)
	if err == nil {
		t.Error("Expected error when saving horse with invalid gender, but got none")
	}

	// Test duplicate email constraint (if exists)
	_, err = user.NewUser(ctx, testDB, "Duplicate User", "constraint@example.com", "stytch-duplicate")
	if err == nil {
		// Note: This test depends on whether you have a unique constraint on email
		// If not, this test will pass and that's okay
		t.Log("No unique email constraint detected")
	}
}