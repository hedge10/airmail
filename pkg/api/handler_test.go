package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hedge10/airmail/constants"
	"github.com/hedge10/airmail/pkg/conf"
	"github.com/stretchr/testify/assert"
)

func delete() {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, constants.MAIL_SERVER_URI+"/messages", nil)
	if err != nil {
		log.Fatal("Cannot create new delete request. " + err.Error())
	}

	_, err = client.Do(req)
	if err != nil {
		log.Fatal("Cannot delete messages. " + err.Error())
		return
	}
}

func get_latest_mail() (*constants.MailpitMessage, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, constants.MAIL_SERVER_URI+"/messages", nil)
	if err != nil {
		log.Fatal(errors.New(err.Error()))
	}

	q := req.URL.Query()
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(errors.New(err.Error()))
	}
	defer resp.Body.Close()

	mr := &constants.MailpitMessagesResponse{}
	body, _ := io.ReadAll(resp.Body)
	json_err := json.Unmarshal(body, mr)

	if json_err != nil {
		return nil, errors.New("Cannot read mailpit api response")
	}

	return &mr.Messages[0], nil
}

func setup() {
	delete()

	os.Setenv("AM_SMTP_PORT", "1025")
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
	setup()

	config, _ := conf.New()

	form := getFormValues()

	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	actual, _ := get_latest_mail()

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "john.doe@example.com", actual.From.Address)
	assert.Equal(t, "Sample email subject", actual.Subject)
}

func TestIncomingMessageHandlerWithFormDataAndRedirectUrl(t *testing.T) {
	setup()

	config, _ := conf.New()

	form := getFormValuesWithRedirect("http://www.google.com")

	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 308, w.Result().StatusCode)
}

func TestIncomingMessageHandlerWithJsonData(t *testing.T) {
	setup()

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
}

func TestIncomingMessageHandlerWithUnknownContentType(t *testing.T) {
	setup()

	config, _ := conf.New()

	r := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	IncomingMessageHandler(config).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestValidation(t *testing.T) {
	type test struct {
		mr   MessageRequest
		want error
	}

	tests := []test{
		{
			mr: MessageRequest{
				SenderAddress: "",
			},
			want: errors.New("sender-address: cannot be blank."),
		},
		{
			mr: MessageRequest{
				SenderAddress: "john.doe@example.com",
				To: []Person{
					{},
				},
			},
			want: errors.New("receiver address: cannot be blank."),
		},
		{
			mr: MessageRequest{
				SenderAddress: "john.doe@example.com",
				To: []Person{
					{
						Name:    "",
						Address: "invalid",
					},
				},
			},
			want: errors.New("receiver address: must be a valid email address."),
		},
		{
			mr: MessageRequest{
				SenderAddress: "john.doe@example.com",
				Redirect:      "lala",
			},
			want: errors.New("_redirect: must be a valid URL."),
		},
		{
			mr: MessageRequest{
				SenderAddress: "john.doe@example.com",
				Redirect:      "http://www.example.com",
				ContentType:   "invalid",
			},
			want: errors.New("_content-type: must be a valid value."),
		},
	}

	for _, tc := range tests {
		got := tc.mr.Validate()
		if !reflect.DeepEqual(tc.want.Error(), got.Error()) {
			t.Fatalf("expected: %v, got: %v", tc.want.Error(), got.Error())
		}
	}
}
