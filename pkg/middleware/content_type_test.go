package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeIsEnforcedToJson(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	EnforceContentType(next).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestContentTypeIsEnforcedToFormUrlEncoded(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	EnforceContentType(next).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestIllegalContentTypeCannotBeHandled(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", "wrong/type")
	w := httptest.NewRecorder()

	EnforceContentType(next).ServeHTTP(w, r)

	assert.Equal(t, 415, w.Result().StatusCode)
}

func TestNotExistingContentTypeCannotBeHandled(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	EnforceContentType(next).ServeHTTP(w, r)

	assert.Equal(t, 400, w.Result().StatusCode)
}
