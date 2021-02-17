package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginHasRedirect(t *testing.T) {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthorizeHandlerFunc)

	form := url.Values{}
	form.Add("username", "peter")

	req, err := http.NewRequest("POST", clientConf.AuthCodeURL("test-state-1"), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.True(t, strings.HasPrefix(rr.Header().Get("Location"), clientConf.RedirectURL))
}
