package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestRandomHash(t *testing.T) {
	str, err := RandomHash(1)
	assert.NoError(t, err)

	h := sha3.New256()
	h.Write([]byte(""))
	assert.NotEqual(t, string(h.Sum(nil)), str)
}
