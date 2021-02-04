package service

import (
	"context"
	"testing"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"beem-auth/internal/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnsurePasswordHash(t *testing.T) {
	a := NewAccountController(db)

	password := "password"
	email := "user@example.com"

	_, err := a.Create(context.Background(), &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)

	user := &database.User{}
	err = db.Get(user, "SELECT email, password FROM users WHERE email=$1", email)
	assert.NoError(t, err)
	assert.NotEqual(t, password, user.Password, "passwords should not be the same, password should be hashed")

	ok, err := util.ComparePasswords(user.Password, password)
	assert.NoError(t, err)
	assert.True(t, ok, "password should match")
}

func TestCreateUserError(t *testing.T) {
	// TODO: change email to user@example.com after middleware transaction rework
	// current architecture doesn't allow us to rollback or access the transaction
	a := NewAccountController(db)

	password := "password"
	email := "user1@example.com"

	_, err := a.Create(context.Background(), &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)

	_, err = a.Create(context.Background(), &pb.AccountCreateRequest{
		Email:    email,
		Password: password,
	})
	assert.Error(t, err)
}
