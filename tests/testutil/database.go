package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/farm"
	"github.com/DevonFarm/sales/horse"
	"github.com/DevonFarm/sales/user"
)

// TestDB wraps database connection with test utilities
type TestDB struct {
	*database.DB
	cleanupFuncs []func() error
}

// NewTestDB creates a test database connection
func NewTestDB(t *testing.T) *TestDB {
	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == "" {
		t.Skip("TEST_DATABASE_URL not set - skipping test")
	}

	db, err := database.NewDBConn(connString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{
		DB:           db,
		cleanupFuncs: make([]func() error, 0),
	}
}

// Cleanup runs all registered cleanup functions
func (db *TestDB) Cleanup() error {
	for _, cleanup := range db.cleanupFuncs {
		if err := cleanup(); err != nil {
			return err
		}
	}
	return nil
}

// CleanupOnExit registers cleanup to run when test exits
func (db *TestDB) CleanupOnExit(t *testing.T) {
	t.Cleanup(func() {
		if err := db.Cleanup(); err != nil {
			t.Errorf("Cleanup failed: %v", err)
		}
		db.Close(context.Background())
	})
}

// TruncateAllTables removes all data from test tables
func (db *TestDB) TruncateAllTables(ctx context.Context) error {
	tables := []string{"horses", "farms", "users"}
	
	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}
	
	return nil
}

// WithTransaction runs a function within a database transaction that is rolled back
func (db *TestDB) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Always rollback for tests

	return fn(tx)
}

// TestFixtures provides test data creation utilities
type TestFixtures struct {
	db *TestDB
}

// NewTestFixtures creates a new fixtures helper
func NewTestFixtures(db *TestDB) *TestFixtures {
	return &TestFixtures{db: db}
}

// CreateTestUser creates a test user and registers cleanup
func (f *TestFixtures) CreateTestUser(ctx context.Context, name, email, stytchID string) (*user.User, error) {
	u, err := user.NewUser(ctx, f.db.DB, name, email, stytchID)
	if err != nil {
		return nil, err
	}

	// Register cleanup
	f.db.cleanupFuncs = append(f.db.cleanupFuncs, func() error {
		_, err := f.db.Exec(ctx, "DELETE FROM users WHERE id = $1", u.ID)
		return err
	})

	return u, nil
}

// CreateTestFarm creates a test farm and registers cleanup
func (f *TestFixtures) CreateTestFarm(ctx context.Context, name string, userID uuid.UUID) (*farm.Farm, error) {
	testFarm := &farm.Farm{Name: name}
	err := testFarm.Save(ctx, f.db.DB, userID.String())
	if err != nil {
		return nil, err
	}

	// Register cleanup
	f.db.cleanupFuncs = append(f.db.cleanupFuncs, func() error {
		_, err := f.db.Exec(ctx, "DELETE FROM farms WHERE id = $1", testFarm.ID)
		return err
	})

	return testFarm, nil
}

// CreateTestHorse creates a test horse and registers cleanup
func (f *TestFixtures) CreateTestHorse(ctx context.Context, h *horse.Horse) (*horse.Horse, error) {
	err := h.Save(ctx, f.db.DB)
	if err != nil {
		return nil, err
	}

	// Register cleanup
	f.db.cleanupFuncs = append(f.db.cleanupFuncs, func() error {
		_, err := f.db.Exec(ctx, "DELETE FROM horses WHERE id = $1", h.ID)
		return err
	})

	return h, nil
}

// CreateTestDataSet creates a complete test dataset with user, farm, and horses
func (f *TestFixtures) CreateTestDataSet(ctx context.Context) (*TestDataSet, error) {
	// Create user
	testUser, err := f.CreateTestUser(ctx, "Test User", "test@example.com", "stytch-123")
	if err != nil {
		return nil, fmt.Errorf("failed to create test user: %w", err)
	}

	// Create farm
	testFarm, err := f.CreateTestFarm(ctx, "Test Farm", testUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create test farm: %w", err)
	}

	// Create horses
	horses := make([]*horse.Horse, 0)
	
	stallion := &horse.Horse{
		Name:        "Thunder",
		Description: "A powerful stallion",
		DateOfBirth: mustParseDate("2018-05-15"),
		Gender:      horse.GenderStallion,
		FarmID:      testFarm.ID,
	}
	stallion, err = f.CreateTestHorse(ctx, stallion)
	if err != nil {
		return nil, fmt.Errorf("failed to create stallion: %w", err)
	}
	horses = append(horses, stallion)

	mare := &horse.Horse{
		Name:        "Beauty",
		Description: "A beautiful mare",
		DateOfBirth: mustParseDate("2019-03-10"),
		Gender:      horse.GenderMare,
		FarmID:      testFarm.ID,
	}
	mare, err = f.CreateTestHorse(ctx, mare)
	if err != nil {
		return nil, fmt.Errorf("failed to create mare: %w", err)
	}
	horses = append(horses, mare)

	gelding := &horse.Horse{
		Name:        "Steady",
		Description: "A calm gelding",
		DateOfBirth: mustParseDate("2017-08-22"),
		Gender:      horse.GenderGelding,
		FarmID:      testFarm.ID,
	}
	gelding, err = f.CreateTestHorse(ctx, gelding)
	if err != nil {
		return nil, fmt.Errorf("failed to create gelding: %w", err)
	}
	horses = append(horses, gelding)

	return &TestDataSet{
		User:   testUser,
		Farm:   testFarm,
		Horses: horses,
	}, nil
}

// TestDataSet represents a complete set of test data
type TestDataSet struct {
	User   *user.User
	Farm   *farm.Farm
	Horses []*horse.Horse
}

// GetStallions returns all stallions from the dataset
func (ds *TestDataSet) GetStallions() []*horse.Horse {
	var stallions []*horse.Horse
	for _, h := range ds.Horses {
		if h.Gender == horse.GenderStallion {
			stallions = append(stallions, h)
		}
	}
	return stallions
}

// GetMares returns all mares from the dataset
func (ds *TestDataSet) GetMares() []*horse.Horse {
	var mares []*horse.Horse
	for _, h := range ds.Horses {
		if h.Gender == horse.GenderMare {
			mares = append(mares, h)
		}
	}
	return mares
}

// GetGeldings returns all geldings from the dataset
func (ds *TestDataSet) GetGeldings() []*horse.Horse {
	var geldings []*horse.Horse
	for _, h := range ds.Horses {
		if h.Gender == horse.GenderGelding {
			geldings = append(geldings, h)
		}
	}
	return geldings
}