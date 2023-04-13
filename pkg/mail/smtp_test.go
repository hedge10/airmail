package mail

import (
	"log"
	"testing"

	"github.com/hedge10/airmail/pkg/conf"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
)

func TestSmtpSendWithoutAuthWithoutAddressNames(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:        2525,
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	m := &Email{
		From: Party{
			Name:  "",
			Email: "john.doe@example.com",
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
		SmtpPort:    2525,
	})

	if te != nil {
		log.Fatal(te)
	}

	err := transfer.Send(m)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(server.Messages()))
	assert.Equal(t, "MAIL FROM:<john.doe@example.com>", server.Messages()[0].MailfromRequest())
	assert.Equal(t, "RCPT TO:<receiver1@example.com>", server.Messages()[0].RcpttoRequestResponse()[0][0])
}
