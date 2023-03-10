package api

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hedge10/airmail/pkg/conf"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *smtpmock.Server {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:        2525,
		LogToStdout:       true,
		LogServerActivity: true,
	})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	os.Setenv("AM_SMTP_PORT", "2525")

	return server
}

func teardown(s *smtpmock.Server) {
	if err := s.Stop(); err != nil {
		fmt.Println(err)
	}
	defer os.Unsetenv("AM_SMTP_PORT")
}

func getFormValuesWithRedirect(redirect string) url.Values {
	form := getFormValues()
	form.Add("_redirect", redirect)

	return form
}

func getFormValues() url.Values {
	form := url.Values{}
	form.Add("sender-address", "john.doe@example.com")
	form.Add("to[0].address", "jane.doe@example.com")
	form.Add("to[1].address", "janice.doe@example.com")
	form.Add("cc[0].address", "paul.doe@example.com")
	form.Add("bcc[0].address", "robert.doe@example.com")
	form.Add("subject", "Sample email subject")
	form.Add("message", "A hilarious message.")

	return form
}

func TestIncomingMessageHandlerWithFormData(t *testing.T) {
	s := setup(t)

	config, _ := conf.New()

	form := getFormValues()

	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	teardown(s)
}

func TestIncomingMessageHandlerWithFormDataAndRedirectUrl(t *testing.T) {
	s := setup(t)

	config, _ := conf.New()

	form := getFormValuesWithRedirect("http://www.google.com")

	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 308, w.Result().StatusCode)

	teardown(s)
}

func TestIncomingMessageHandlerWithJsonData(t *testing.T) {
	s := setup(t)

	json := []byte(`{
		"sender-address": "john.doe@example.com",
		"to": [
			{
				"address": "jane.doe@example.com"
			},
			{
				"address": "janice.doe@example.com"
			}
		],
		"cc": [
			{
				"address": "paul.doe@example.com"
			},
			{
				"address": "paula.doe@example.com"
			}
		],
		"bcc": [
			{
				"address": "robert.doe@example.com"
			},
			{
				"address": "roberta.doe@example.com"
			}
		],
		"subject": "Sample email subject",
		"message": "A hilarious message."
	}`)

	config, _ := conf.New()

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(json))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	teardown(s)
}

func TestIncomingMessageHandlerWithUnknownContentType(t *testing.T) {
	s := setup(t)

	config, _ := conf.New()

	r := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	teardown(s)
}
