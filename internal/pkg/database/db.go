package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect(host string, port string, user string, password string, dbName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}


	return db, nil
}

func dbAccessError(err error) error {
	return fmt.Errorf("accessing DB: %w", err)
}
