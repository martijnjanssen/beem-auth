package middleware

import (
	"beem-auth/internal/pkg/database"
	"context"
	"database/sql"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTransactionInterceptor(t *testing.T) {
	td, iDb := database.StartTestPostgreSQL()
	assert.NoError(t, database.ApplyMigrations(iDb))
	// Teardown of started testing database
	defer td()

	interceptor := NewTransactionInterceptor(iDb)

	email := "user@example.com"
	handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
		tx := GetContextTx(ctx)
		err := database.UserAdd(ctx, tx, email, "hashedPassword")

		return nil, err
	}

	_, err := interceptor(context.Background(), nil, nil, handlerFunc)
	assert.NoError(t, err)

	_, err = database.UserGetOnEmail(context.Background(), iDb, email)
	assert.NoError(t, err, "user should be found, transaction would be committed")
}

func TestTransactionInterceptorRollback(t *testing.T) {
	td, iDb := database.StartTestPostgreSQL()
	assert.NoError(t, database.ApplyMigrations(iDb))
	// Teardown of started testing database
	defer td()

	interceptor := NewTransactionInterceptor(iDb)

	email := "user@example.com"
	handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
		tx := GetContextTx(ctx)
		err := database.UserAdd(ctx, tx, email, "hashedPassword")
		assert.NoError(t, err)

		return nil, status.Errorf(codes.Internal, "this is an error to trigger a rollback")
	}

	_, err := interceptor(context.Background(), nil, nil, handlerFunc)
	assert.Error(t, err, "the error from the handerfunc should be passed here")

	_, err = database.UserGetOnEmail(context.Background(), iDb, email)
	assert.True(t, errors.Is(err, sql.ErrNoRows), "user should not be found, transaction would be rolled back")
}

func TestTransactionInterceptorTxStartFailed(t *testing.T) {
	interceptor := NewTransactionInterceptor(closedDb)

	handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	_, err := interceptor(context.Background(), nil, nil, handlerFunc)
	assert.Error(t, err, "should be unable to start the transaction")
}

func TestSetGetContextTransaction(t *testing.T) {
	ctx := context.Background()
	tx := db.MustBegin()

	txCtx := SetContextTx(ctx, tx)
	ctxTx := GetContextTx(txCtx)

	_, err := ctxTx.Exec("SELECT * FROM users")
	assert.NoError(t, err)

	err = tx.Rollback()
	assert.NoError(t, err)
}
