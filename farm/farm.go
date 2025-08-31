package farm

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DevonFarm/sales/database"
)

type Farm struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func NewFarm(ctx context.Context, name string, db *database.DB) (*Farm, error) {
	farm := &Farm{
		Name: name,
	}
	if err := farm.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to save farm: %w", err)
	}
	return farm, nil
}

func (f *Farm) Save(ctx context.Context, db *database.DB) error {
	row := db.QueryRow(
		ctx,
		`INSERT INTO farms (name) VALUES ($1) RETURNING id`,
		f.Name,
	)
	return row.Scan(&f.ID)
}

type User struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	Email  string    `db:"email"`
	FarmID uuid.UUID `db:"farm_id"`
}

func (f *Farm) NewUser(ctx context.Context, name, email string, db *database.DB) (*User, error) {
	user := &User{
		Name:   name,
		Email:  email,
		FarmID: f.ID,
	}
	if err := user.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	return user, nil
}

func (u *User) Save(ctx context.Context, db *database.DB) error {
	row := db.QueryRow(
		ctx,
		`INSERT INTO users (name, email, farm_id) VALUES ($1, $2, $3) RETURNING id`,
		u.Name,
		u.Email,
		u.FarmID,
	)
	return row.Scan(&u.ID)
}

func GetFarmByUser(userID uuid.UUID) (*Farm, error) {
	// TODO: implement database query to get farm by user ID

	return nil, nil
}
