package service

import (
	"context"
	"testing"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"beem-auth/internal/pkg/middleware"
	"beem-auth/internal/pkg/util/email"
	"beem-auth/internal/pkg/util/hash"

	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mail = "user@example.com"
var password = "password"

func TestCreateEnsurePasswordHash(t *testing.T) {
	a, e, ctx, tx, rb := accountControllerHelper(db)
	defer rb()

	e.On("SendEmail", mock.Anything).Return(nil)

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    mail,
		Password: password,
	})
	assert.NoError(t, err)

	user := &database.User{}
	err = tx.Get(user, "SELECT email, password FROM users WHERE email=$1", mail)
	assert.NoError(t, err)
	assert.NotEqual(t, password, user.Password, "passwords should not be the same, password should be hashed")

	ok, err := hash.ComparePasswords(user.Password, password)
	assert.NoError(t, err)
	assert.True(t, ok, "password should match")

	e.AssertCalled(t, "SendEmail", mock.MatchedBy(func(e email.Email) bool { return e.Recipient == mail }))
}

func TestCreateUserError(t *testing.T) {
	a, e, ctx, _, rb := accountControllerHelper(db)
	defer rb()

	e.On("SendEmail", mock.Anything).Return(nil)

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    mail,
		Password: password,
	})
	assert.NoError(t, err)

	_, err = a.Create(ctx, &pb.AccountCreateRequest{
		Email:    mail,
		Password: password,
	})
	assert.Error(t, err)
}

func TestCreateUserEmailError(t *testing.T) {
	a, e, ctx, _, rb := accountControllerHelper(db)
	defer rb()

	e.On("SendEmail", mock.Anything).Return(fmt.Errorf("error for testing"))

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    mail,
		Password: password,
	})
	assert.Error(t, err)
}

func TestCreateUserChallengeCreateError(t *testing.T) {
	a, e, ctx, tx, rb := accountControllerHelper(db)
	defer rb()

	_ = tx.MustExec("DROP TABLE challenges")

	e.On("SendEmail", mock.Anything).Return(nil)

	_, err := a.Create(ctx, &pb.AccountCreateRequest{
		Email:    mail,
		Password: password,
	})
	assert.Error(t, err)
}

func accountControllerHelper(db *sqlx.DB) (pb.AccountServiceServer, *email.EmailMock, context.Context, *sqlx.Tx, func() error) {
	e := email.NewEmailMock()

	a := NewAccountController(e)
	tx := db.MustBegin()
	rb := tx.Rollback
	ctx := middleware.SetContextTx(context.Background(), tx)

	return a, e, ctx, tx, rb
}
