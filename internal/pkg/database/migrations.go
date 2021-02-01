package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var migrations = []migration{
	{
		name: "1_users_create_table",
		fn: func(db Queryer) error {
			_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (email text, password text)")
			return err
		},
	}, {
		name: "2_users_add_email_uniqueness",
		fn: func(db Queryer) error {
			_, err := db.Exec("ALTER TABLE users ADD UNIQUE (email);")
			return err
		},
	},
}

func ApplyMigrations(db *sqlx.DB) error {
	err := migrate(db, migrations)
	if err != nil {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return nil
}
