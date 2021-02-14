package web

import (
	"log"
	"os"

	"testing"

	"beem-auth/internal/pkg/database"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var closedDb *sqlx.DB

// Initializes the database for the tests run in this package
func TestMain(m *testing.M) {
	// Functions for teardown of started docker containers
	var td, closedTd func()

	td, db = database.StartTestPostgreSQL()
	closedTd, closedDb = database.StartTestPostgreSQL()
	if err := closedDb.Close(); err != nil {
		log.Fatalf("unable to close database: %s", err)
	}
	closedTd()

	if err := database.ApplyMigrations(db); err != nil {
		log.Fatalf("migration failed: %s", err)
	}

	code := m.Run()
	td()
	os.Exit(code)
}
