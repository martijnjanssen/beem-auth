package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// From https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

func HashAndSalt(pwd string) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	bytePlain := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePlain, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %w", err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}
func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to compare passwords: %w", err)
	}

	return true, nil
}
