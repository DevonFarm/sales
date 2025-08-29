package farm

import "github.com/google/uuid"

type Farm struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type User struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	Email  string    `db:"email"`
	FarmID uuid.UUID `db:"farm_id"`
}

func GetFarmByUser(userID uuid.UUID) (*Farm, error) {
	// TODO: implement database query to get farm by user ID
	return nil, nil
}
