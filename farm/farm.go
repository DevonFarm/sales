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

// TODO: need a /farm/new route to create a farm
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

func GetUserByStytchID(ctx context.Context, db *database.DB, stytchID string) (*User, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, name, email, farm_id, stytch_id FROM users WHERE stytch_id = $1`,
		stytchID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by stytch_id: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by stytch_id: %w", err)
	}
	return &user, nil
}

func (u *User) Save(ctx context.Context, db *database.DB) error {
	found := false
	if u.ID != uuid.Nil {
		found = true
	} else {
		foundUser, err := GetUserByStytchID(ctx, db, u.StytchID)
		if err != nil {
			return fmt.Errorf("failed to check existing user: %w", err)
		}
		if foundUser != nil {
			found = true
			u.ID = foundUser.ID
		} else {
			row := db.QueryRow(
				ctx,
				`INSERT INTO users (name, email, stytch_id, farm_id) VALUES ($1, $2, $3, NULLIF($4, $5)) RETURNING id`,
				u.Name,     // $1
				u.Email,    // $2
				u.StytchID, // $3
				u.FarmID,   // $4
				uuid.Nil,   // $5
			)
			return row.Scan(&u.ID)
		}
	}
	if found {
		// Update existing user
		_, err := db.Exec(
			ctx,
			`UPDATE users SET name = $1, email = $2, farm_id = NULLIF($3, $4) WHERE id = $5`,
			u.Name,   // $1
			u.Email,  // $2
			u.FarmID, // $3
			uuid.Nil, // $4
			u.ID,     // $5
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}
	return nil
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
