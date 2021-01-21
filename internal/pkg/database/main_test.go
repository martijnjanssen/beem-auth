package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	"testing"
)

var db *sqlx.DB

// Initializes the database for the tests run in this package
func TestMain(m *testing.M) {
	var td func()
	td, db = StartTestPostgreSQL()
	code := m.Run()
	td()
	os.Exit(code)
}
