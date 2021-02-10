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
	}, {
		name: "3_users_add_validity",
		fn: func(db Queryer) error {
			_, err := db.Exec("ALTER TABLE users ADD COLUMN valid boolean DEFAULT false;")
			return err
		},
	}, {
		name: "4_users_add_id",
		fn: func(db Queryer) error {
			_, err := db.Exec("ALTER TABLE users ADD COLUMN id uuid DEFAULT uuid_generate_v4();")
			return err
		},
	}, {
		name: "5_challenge_create_table",
		fn: func(db Queryer) error {
			_, err := db.Exec("CREATE TABLE IF NOT EXISTS challenges (user_id uuid, key text UNIQUE)")
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
