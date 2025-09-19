package horse

import (
	"context"
	"fmt"
	"time"

	"github.com/DevonFarm/sales/database"

	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

type Horse struct {
	ID          uuid.UUID `db:"id" form:"id"`
	Name        string    `db:"name" form:"name"`
	Description string    `db:"description" form:"description"`
	Images      []*Image
	DateOfBirth time.Time `db:"date_of_birth" form:"-"`
	Gender      Gender    `db:"gender" form:"gender"`
	FarmID      uuid.UUID `db:"farm_id" form:"-"`
}

func (h *Horse) Age() int {
	return int(time.Since(h.DateOfBirth).Hours() / 24 / 365)
}

func (h *Horse) HTMLPath() string {
	return fmt.Sprintf("%s.html", strcase.ToSnake(h.Name))
}

func (h *Horse) NewImage(full, thumbnail, alt string) {
	prefix := fmt.Sprintf("assets/images/horses/%s", h.Name)
	img := &Image{
		Full:      fmt.Sprintf("%s/%s", prefix, full),
		Thumbnail: fmt.Sprintf("%s/%s", prefix, thumbnail),
		Alt:       alt,
	}
	if h.Images == nil {
		h.Images = make([]*Image, 0)
	}
	h.Images = append(h.Images, img)
}

func (h *Horse) Save(ctx context.Context, db *database.DB) error {
	// Validate horse data
	if h.Gender.IsInvalid() {
		return fmt.Errorf("invalid horse gender: %d", h.Gender)
	}
	row := db.QueryRow(
		ctx,
		`INSERT INTO horses (name, description, date_of_birth, gender, farm_id) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		h.Name,        // $1
		h.Description, // $2
		h.DateOfBirth, // $3
		h.Gender,      // $4
		h.FarmID,      // $5
	)
	if err := row.Scan(&h.ID); err != nil {
		return fmt.Errorf("failed to scan horse id: %w", err)
	}
	return nil
}

func GetHorsesByFarmID(ctx context.Context, db *database.DB, farmID uuid.UUID) ([]*Horse, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, name, description, date_of_birth, gender FROM horses WHERE farm_id = $1 ORDER BY name`,
		farmID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query horses: %w", err)
	}
	defer rows.Close()

	var horses []*Horse
	for rows.Next() {
		var h Horse
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.DateOfBirth, &h.Gender); err != nil {
			return nil, fmt.Errorf("failed to scan horse: %w", err)
		}
		horses = append(horses, &h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return horses, nil
}

type DashboardStats struct {
	TotalHorses int
	Stallions   int
	Mares       int
	Geldings    int
}

func GetDashboardStats(ctx context.Context, db *database.DB, farmID uuid.UUID) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Get total count
	row := db.QueryRow(ctx, `SELECT COUNT(*) FROM horses WHERE farm_id = $1`, farmID)
	if err := row.Scan(&stats.TotalHorses); err != nil {
		return nil, fmt.Errorf("failed to get total horses: %w", err)
	}

	// Get gender breakdown
	rows, err := db.Query(ctx, `SELECT gender, COUNT(*) FROM horses WHERE farm_id = $1 GROUP BY gender`, farmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var gender Gender
		var count int
		if err := rows.Scan(&gender, &count); err != nil {
			return nil, fmt.Errorf("failed to scan gender stats: %w", err)
		}
		switch gender {
		case GenderStallion:
			stats.Stallions = count
		case GenderMare:
			stats.Mares = count
		case GenderGelding:
			stats.Geldings = count
		}
	}

	return stats, nil
}

type Image struct {
	Full      string
	Alt       string
	Thumbnail string
}
