package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	*pgx.Conn
}

func NewDBConn(connString string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &DB{Conn: conn}, nil
}
