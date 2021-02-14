package web

import (
	"testing"

	"beem-auth/internal/pkg/database"

	"context"
	"net/http"
	"net/http/httptest"

	"errors"
	"github.com/stretchr/testify/assert"
)

func TestWrapTransactionCommit(t *testing.T) {
	td, iDb := database.StartTestPostgreSQL()
	assert.NoError(t, database.ApplyMigrations(iDb))
	// Teardown of started testing database
	defer td()

	// Set the database for the transactionWrapper
	SetDB(iDb)

	var handler http.HandlerFunc = wrapTransaction(func(ctx context.Context, db database.Queryer, w http.ResponseWriter, req *http.Request) (func(), error) {
		_, err := database.UserAdd(ctx, db, "user@example.com", "password")
		assert.NoError(t, err)

		return func() {
			writeResponse(w, http.StatusOK, "user added")
		}, nil
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	_, err = database.UserGetOnEmail(context.Background(), iDb, "user@example.com")
	assert.NoError(t, err)
}

func TestWrapTransactionRollback(t *testing.T) {
	td, iDb := database.StartTestPostgreSQL()
	assert.NoError(t, database.ApplyMigrations(iDb))
	// Teardown of started testing database
	defer td()

	// Set the database for the transactionWrapper
	SetDB(iDb)

	var handler http.HandlerFunc = wrapTransaction(func(ctx context.Context, db database.Queryer, w http.ResponseWriter, req *http.Request) (func(), error) {
		_, err := database.UserAdd(ctx, db, "user@example.com", "password")
		assert.NoError(t, err)

		writeResponse(w, http.StatusBadRequest, "invalid input")

		return nilFn, errors.New("this is a testing error")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	_, err = database.UserGetOnEmail(context.Background(), iDb, "user@example.com")
	assert.Error(t, err)
}

func TestWrapTransactionUnableToStartTx(t *testing.T) {
	// Set the database for the transactionWrapper
	SetDB(closedDb)

	var handler http.HandlerFunc = wrapTransaction(func(ctx context.Context, db database.Queryer, w http.ResponseWriter, req *http.Request) (func(), error) {
		return nilFn, nil
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
