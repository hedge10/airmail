package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setupServer() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = CreateClient(BaseUri(server.URL))

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := os.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestGrecaptchaWithSuccessfulResponse(t *testing.T) {
	teardown := setupServer()
	defer teardown()

	mux.HandleFunc("/recaptcha/api/siteverify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("200-with-success.json"))
	})

	err := client.ValidateGrecaptcha("secret", "recaptcha-response", "127.0.0.1")

	assert.Nil(t, err)
}

func TestGrecaptchaWithErrorResponse(t *testing.T) {
	teardown := setupServer()
	defer teardown()

	mux.HandleFunc("/recaptcha/api/siteverify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("200-with-error.json"))
	})

	err := client.ValidateGrecaptcha("secret", "recaptcha-response", "127.0.0.1")

	assert.NotNil(t, err)
	assert.Equal(t, "unsuccessful recaptcha verify request", err.Error())
}
