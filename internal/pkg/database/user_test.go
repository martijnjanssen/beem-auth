package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreateGet(t *testing.T) {
	tx := db.MustBegin()

	assert.NoError(t, UserAdd(context.Background(), tx, "user1@example.com", "password"))
	assert.NoError(t, UserAdd(context.Background(), tx, "user2@example.com", "password"))

	user, err := UserGetOnEmail(context.Background(), tx, "user1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, "user1@example.com")

	assert.NoError(t, tx.Rollback())
}

func TestUserCreateDuplicate(t *testing.T) {
	tx := db.MustBegin()

	assert.NoError(t, UserAdd(context.Background(), tx, "user1@example.com", "password"))
	assert.Error(t, UserAdd(context.Background(), tx, "user1@example.com", "password"))

	assert.NoError(t, tx.Rollback())
}

func TestUserCreateErrorAccess(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	err := UserAdd(context.Background(), tx, "user@example.com", "password")

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}

func TestUserGetErrorNotFound(t *testing.T) {
	_, err := UserGetOnEmail(context.Background(), db, "notfound@example.com")
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}
