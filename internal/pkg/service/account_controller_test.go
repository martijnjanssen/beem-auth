package service

import (
	"context"
	"testing"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"beem-auth/internal/pkg/middleware"
	"beem-auth/internal/pkg/util"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnsurePasswordHash(t *testing.T) {
	a, ctx, tx, rb := accountControllerHelper(db)
	defer rb()

	password := "password"
	email := "user@example.com"

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)

	user := &database.User{}
	err = tx.Get(user, "SELECT email, password FROM users WHERE email=$1", email)
	assert.NoError(t, err)
	assert.NotEqual(t, password, user.Password, "passwords should not be the same, password should be hashed")

	ok, err := util.ComparePasswords(user.Password, password)
	assert.NoError(t, err)
	assert.True(t, ok, "password should match")
}

func TestCreateUserError(t *testing.T) {
	a, ctx, _, rb := accountControllerHelper(db)
	defer rb()

	password := "password"
	email := "user@example.com"

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)

	_, err = a.Create(ctx, &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.Error(t, err)
}

func accountControllerHelper(db *sqlx.DB) (pb.AccountServiceServer, context.Context, *sqlx.Tx, func() error) {
	a := NewAccountController()
	tx := db.MustBegin()
	rb := tx.Rollback
	ctx := middleware.SetContextTx(context.Background(), tx)

	return a, ctx, tx, rb

}
