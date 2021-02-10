package database

import (
	"context"
	"testing"

	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChallengeCreate(t *testing.T) {
	tx := db.MustBegin()

	userId := uuid.New()

	str, err := ChallengeCreate(context.Background(), tx, userId)
	assert.NoError(t, err)

	challenges := []Challenge{}
	assert.NoError(t, tx.Select(&challenges, "SELECT * FROM challenges"))
	assert.Len(t, challenges, 1, "expected only one challenge to be found")
	assert.Equal(t, userId, challenges[0].UserId)
	assert.Equal(t, str, challenges[0].Key)

	assert.NoError(t, tx.Rollback())
}

func TestChallengeAddDBAccessError(t *testing.T) {
	tx := db.MustBegin()
	assert.NoError(t, tx.Rollback())

	_, err := ChallengeCreate(context.Background(), tx, uuid.New())

	assert.True(t, errors.Is(err, sql.ErrTxDone))
}
