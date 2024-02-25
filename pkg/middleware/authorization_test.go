package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectAuthTokenGiven(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer 123")
	w := httptest.NewRecorder()

	auth := NewToken("123")
	auth.Validate(next).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestNoAuthTokenGiven(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	auth := NewToken("Bearer 123")
	auth.Validate(next).ServeHTTP(w, r)

	assert.Equal(t, 400, w.Result().StatusCode)
}

func TestWrongAuthTokenGiven(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer wrong-token")
	w := httptest.NewRecorder()

	auth := NewToken("Bearer 123")
	auth.Validate(next).ServeHTTP(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)
}
