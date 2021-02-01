package database

import (
	"database/sql"
	"errors"
	"log"

	"fmt"
	"github.com/jmoiron/sqlx"
)

// TODO: pick correct type for timestamp (with/without timezone)
var migrationSchema = `
  CREATE TABLE IF NOT EXISTS migrations (
    name text UNIQUE,
    timestamp timestamp
  )
`

type migrationFunc func(Queryer) error

type migration struct {
	name string
	fn   migrationFunc
}

type migrationLog struct {
	Name      string `db:"name"`
	Timestamp string `db:"timestamp"`
}

func migrate(db *sqlx.DB, migrations []migration) error {
	err := createMigrationTable(db)
	if err != nil {
		return fmt.Errorf("unable to create migration table: %w", err)
	}

	for _, m := range migrations {
		err = applyMigration(db, m.name, m.fn)
		if err != nil {
			return fmt.Errorf("unable to apply migration '%s': %w", m.name, err)
		}
	}

	return nil
}

func createMigrationTable(db *sqlx.DB) error {
	_, err := db.Exec(migrationSchema)
	if err != nil {
		return dbAccessError(err)
	}

	return nil
}

// migrate applies a migration and stops execution completely when applying a migration fails
func applyMigration(db *sqlx.DB, name string, fn func(Queryer) error) error {
	tx := db.MustBegin()

	// Check if the migration is already applied
	applied, err := checkMigrationAlreadyApplied(tx, name)
	if err != nil {
		rollbackMigrationTx(tx)
		return fmt.Errorf("error while checking for applied migrations: %w", err)
	}

	// If migration is already applied, skip it
	if applied {
		rollbackMigrationTx(tx)
		return nil
	}

	// Keep track of which migrations have been applied
	err = addMigrationLog(tx, name)
	if err != nil {
		rollbackMigrationTx(tx)
		return fmt.Errorf("error while adding migration tracker: %w", err)
	}

	// Execute the migration itself
	err = fn(tx)
	if err != nil {
		rollbackMigrationTx(tx)
		return fmt.Errorf("error while applying migration: %w", err)
	}

	// Commit migration to the database
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while committing transaction: %w", err)
	}

	return nil
}

func checkMigrationAlreadyApplied(db Queryer, name string) (bool, error) {
	migration := &migrationLog{}
	err := db.Get(migration, "SELECT name FROM migrations WHERE name=$1", name)
	if err != nil {
		// If there are no rows, the migration is not yet applied
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		// If we reach this there was a normal error
		return true, err
	}

	// We got the migration we asked for, it is applied
	return true, nil
}

func addMigrationLog(db Queryer, name string) error {
	_, err := db.Exec("INSERT INTO migrations (name, timestamp) VALUES ($1, now())", name)
	if err != nil {
		return err
	}

	return nil
}

// rollbackTx is a shorter method for rolling back while in a migration
func rollbackMigrationTx(tx *sqlx.Tx) {
	if rollErr := tx.Rollback(); rollErr != nil {
		log.Printf("error while rolling back failed migration: %s", rollErr)
	}
}
