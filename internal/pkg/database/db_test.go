package database

import (
	"database/sql"
	"errors"
	"testing"

	"os"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestDBAccessError(t *testing.T) {
	err := dbAccessError(sql.ErrNoRows)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.False(t, errors.Is(err, sql.ErrTxDone))
}

func TestMain(m *testing.M) {
	var td func()
	td, db = StartTestPostgreSQL()
	code := m.Run()
	td()
	os.Exit(code)
}
