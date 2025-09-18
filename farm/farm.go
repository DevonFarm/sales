package farm

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/user"
)

type Farm struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func NewFarm(ctx context.Context, name string, db *database.DB, userID string) (*Farm, error) {
	farm := &Farm{
		Name: name,
	}
	if err := farm.Save(ctx, db, userID); err != nil {
		return nil, fmt.Errorf("failed to save farm: %w", err)
	}
	return farm, nil
}

func (f *Farm) Save(ctx context.Context, db *database.DB, userID string) error {
	if f.ID == uuid.Nil {
		u, err := user.GetUser(ctx, db, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if u == nil {
			return fmt.Errorf("user not found with ID: %s", userID)
		}
		if u.FarmID != uuid.Nil {
			f.ID = u.FarmID
		} else {
			row := db.QueryRow(
				ctx,
				`INSERT INTO farms (name) VALUES ($1) RETURNING id`,
				f.Name,
			)
			if err := row.Scan(&f.ID); err != nil {
				return fmt.Errorf("failed to insert farm: %w", err)
			}
			// Associate the farm with the user
			_, err := db.Exec(
				ctx,
				`UPDATE users SET farm_id = $1 WHERE id = $2`,
				f.ID,
				userID,
			)
			if err != nil {
				return fmt.Errorf("failed to associate farm with user: %w", err)
			}
			return nil
		}
	}
	// Update existing farm
	_, err := db.Exec(
		ctx,
		`UPDATE farms SET name = $1 WHERE id = $2`,
		f.Name,
		f.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update farm: %w", err)
	}
	return nil
}

func GetFarm(ctx context.Context, db *database.DB, farmID string) (*Farm, error) {
	var farm Farm
	row := db.QueryRow(
		ctx,
		`SELECT id, name FROM farms WHERE id = $1`,
		farmID,
	)
	if err := row.Scan(&farm.ID, &farm.Name); err != nil {
		return nil, fmt.Errorf("failed to get farm: %w", err)
	}
	return &farm, nil
}
