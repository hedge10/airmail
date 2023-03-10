package mail

import (
	"fmt"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendWithoutAuthWithoutAddressNames(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:        2525,
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	print(server.Start().Error())

	m := Email{
		Connection: Connection{
			Address: fmt.Sprintf("127.0.0.1:%d", server.PortNumber()),
		},
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

	err := m.SendWithoutAuth()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(server.Messages()))
	assert.Equal(t, "MAIL FROM:<john.doe@example.com>", server.Messages()[0].MailfromRequest())
	assert.Equal(t, "RCPT TO:<receiver1@example.com>", server.Messages()[0].RcpttoRequestResponse()[0][0])
}
