package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserCreateGet(t *testing.T) {
	tx := db.MustBegin()

	userId1, err := UserAdd(context.Background(), tx, "user1@example.com", "password")
	assert.NoError(t, err)
	_, err = UserAdd(context.Background(), tx, "user2@example.com", "password")
	assert.NoError(t, err)
	assert.NotEqual(t, userId1, uuid.Nil)

	user, err := UserGetOnEmail(context.Background(), tx, "user1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, "user1@example.com")
	assert.Equal(t, userId1, user.Id)

	assert.NoError(t, tx.Rollback())
}

func TestUserCreateDuplicate(t *testing.T) {
	tx := db.MustBegin()

	_, err := UserAdd(context.Background(), tx, "user1@example.com", "password")
	assert.NoError(t, err)
	_, err = UserAdd(context.Background(), tx, "user1@example.com", "password")
	assert.Error(t, err)

	assert.NoError(t, tx.Rollback())
}

func TestUserCreateErrorAccess(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	_, err := UserAdd(context.Background(), tx, "user@example.com", "password")

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}

func TestUserGetErrorNotFound(t *testing.T) {
	_, err := UserGetOnEmail(context.Background(), db, "notfound@example.com")
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestUserSetValid(t *testing.T) {
	tx := db.MustBegin()

	userId, err := UserAdd(context.Background(), tx, "user@example.com", "password")
	assert.NoError(t, err)

	err = UserSetValid(context.Background(), tx, userId)
	assert.NoError(t, err)

	user, err := UserGetOnEmail(context.Background(), tx, "user@example.com")
	assert.NoError(t, err)
	assert.True(t, user.Valid)

	assert.NoError(t, tx.Rollback())
}

func TestUserSetValidErrorAccess(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	err := UserSetValid(context.Background(), tx, uuid.New())

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}

func TestUserSetValidErrorNotFound(t *testing.T) {
	err := UserSetValid(context.Background(), db, uuid.New())
	assert.Error(t, err)
}
