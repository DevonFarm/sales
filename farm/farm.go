package farm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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
	ID       uuid.UUID `db:"id" form:"-"`
	Name     string    `db:"name" form:"name"`
	Email    string    `db:"email" form:"email"`
	FarmID   uuid.UUID `db:"farm_id" form:"-"`
	StytchID string    `db:"stytch_id" form:"-"`
}

func NewUser(ctx context.Context, db *database.DB, name, email, stytchID string) (*User, error) {
	user := &User{
		Name:     name,
		Email:    email,
		StytchID: stytchID,
	}
	if err := user.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	return user, nil
}

func (u *User) Save(ctx context.Context, db *database.DB) error {
	row := db.QueryRow(
		ctx,
		`INSERT INTO users (name, email, stytch_id) VALUES ($1, $2, $3) RETURNING id`,
		u.Name,
		u.Email,
		u.StytchID,
	)
	return row.Scan(&u.ID)
}

func GetFarmByUser(ctx context.Context, db *database.DB, userID uuid.UUID) (*Farm, error) {
	rows, err := db.Query(
		ctx,
		`SELECT f.id, f.name FROM farms f JOIN users u ON f.id = u.farm_id WHERE u.id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query farm by user: %w", err)
	}
	farm, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Farm])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get farm by user: %w", err)
	}
	return &farm, nil
}
