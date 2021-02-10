package random

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Generates a random string with a length len
// which is hashed with sha3-256.
func RandomHash(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error while reading random string: %w", err)
	}

	h := sha3.New256()
	h.Write(b)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash, nil
}
