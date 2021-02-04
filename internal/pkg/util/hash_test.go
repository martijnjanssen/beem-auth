package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordFlow(t *testing.T) {
	hashedPass, err := HashAndSalt("password")
	assert.NoError(t, err, "expected no error after hashing password")

	ok, err := ComparePasswords(hashedPass, "password")
	assert.NoError(t, err, "expected no error after comparing")

	assert.True(t, ok, "expected passwords to be equal")
}

func TestWrongPasswordFlow(t *testing.T) {
	hashedPass, err := HashAndSalt("password")
	assert.NoError(t, err, "expected no error after hashing password")

	ok, err := ComparePasswords(hashedPass, "wrongpassword")
	assert.NoError(t, err, "expected no error after comparing")

	assert.False(t, ok, "expected passwords to be unequal")
}

func TestInvalidHash(t *testing.T) {
	ok, err := ComparePasswords("invalidHash", "invalidHash")
	assert.Error(t, err, "expected an error")
	assert.True(t, errors.Is(err, bcrypt.ErrHashTooShort), "expected the hash to be invalid")
	assert.False(t, ok, "ok should be false, password check should fail")
}
