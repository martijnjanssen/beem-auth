package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	SetDB(db)

	close := ListenHTTP()

	defer close()

	res, err := http.Get("http://localhost:8081/challenge/complete")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
