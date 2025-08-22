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
	ID          uuid.UUID
	Name        string
	Description string
	Images      []*Image
	DateOfBirth time.Time
	gender      Gender
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

func NewHorse(name string, desc string, dob time.Time, g Gender) *Horse {
	return &Horse{
		Name:        name,
		Description: desc,
		DateOfBirth: dob,
		gender:      g,
	}
}

func (h *Horse) Save(ctx context.Context, db *database.DB) error {
	row := db.QueryRow(
		ctx,
		`INSERT INTO horses (name, description, date_of_birth, gender) 
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		h.Name, h.Description, h.DateOfBirth, h.gender,
	)
	if err := row.Scan(&h.ID); err != nil {
		return fmt.Errorf("failed to scan horse id: %w", err)
	}
	return nil
}

type Image struct {
	Full      string
	Alt       string
	Thumbnail string
}
