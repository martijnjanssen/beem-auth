package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBAccessError(t *testing.T) {
	err := dbAccessError(sql.ErrNoRows)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.False(t, errors.Is(err, sql.ErrTxDone))
}
