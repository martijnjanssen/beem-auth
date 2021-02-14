package web

import (
	"beem-auth/internal/pkg/database"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChallengeCompleteHandler(t *testing.T) {
	tx := db.MustBegin()

	userId, err := database.UserAdd(context.Background(), tx, "user@example.com", "")
	assert.NoError(t, err)

	key, err := database.ChallengeCreate(context.Background(), tx, userId)
	assert.NoError(t, err)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/challenge/complete?key=%s", key), nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	callback, err := challengeCompleteHandler(context.Background(), tx, rr, req)
	assert.NoError(t, err)
	// Trigger writing response with the callback
	callback()

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.NoError(t, tx.Rollback())
}

func TestChallengeCompleteHandlerNoKeyError(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/challenge/complete", nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	_, err = challengeCompleteHandler(context.Background(), closedDb, rr, req)
	assert.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestChallengeCompleteHandlerNoChallengeError(t *testing.T) {
	tx := db.MustBegin()

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/challenge/complete?key=randomkey"), nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	_, err = challengeCompleteHandler(context.Background(), tx, rr, req)
	assert.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.NoError(t, tx.Rollback())
}

func TestChallengeCompleteHandlerNoUserError(t *testing.T) {
	tx := db.MustBegin()

	key, err := database.ChallengeCreate(context.Background(), tx, uuid.New())
	assert.NoError(t, err)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/challenge/complete?key=%s", key), nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	_, err = challengeCompleteHandler(context.Background(), tx, rr, req)
	assert.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.NoError(t, tx.Rollback())
}
