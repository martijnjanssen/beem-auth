package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMigrate(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	err := userMigrate(tx)

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}

func TestUserCreateGet(t *testing.T) {
	tx := db.MustBegin()

	assert.NoError(t, UserAdd(tx, "user1@example.com", "password"))
	assert.NoError(t, UserAdd(tx, "user2@example.com", "password"))

	user, err := UserGetOnEmail(tx, "user1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, "user1@example.com")

	assert.NoError(t, tx.Rollback())
}

func TestUserCreateErrorAccess(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	err := UserAdd(tx, "user@example.com", "password")

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}

func TestUserGetErrorNotFound(t *testing.T) {
	_, err := UserGetOnEmail(db, "notfound@example.com")
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}
