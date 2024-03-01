package mail

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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

func TestSmtpSendWithoutAuthWithoutAddressNames(t *testing.T) {
	setup()

	m := &Email{
		From: Party{
			Name:  "",
			Email: "johnny.boy@example.com",
		},
		To: []Party{
			{
				Name:  "",
				Email: "receiver1@example.com",
			},
		},
		Cc: []Party{
			{
				Name:  "",
				Email: "cc@example.com",
			},
		},
		Subject: "Testing subject",
		Message: "A fantastic email.",
	}

	transfer, te := CreateTransfer(&conf.Config{
		MailService: "smtp",
		SmtpAuth:    "none",
		SmtpHost:    "localhost",
		SmtpPort:    1025,
	})

	if te != nil {
		log.Fatal(te)
	}

	err := transfer.Send(m)
	actual, _ := get_latest_mail()

	assert.Nil(t, err)

	assert.Equal(t, "johnny.boy@example.com", actual.From.Address)
	assert.Equal(t, "receiver1@example.com", actual.To[0].Address)
	assert.Equal(t, "cc@example.com", actual.Cc[0].Address)
	assert.Equal(t, "Testing subject", actual.Subject)
}
