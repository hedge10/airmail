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

	from := party{
		Name:  "",
		Email: "john.doe@example.com",
	}
	to := []party{
		{
			Name:  "",
			Email: "receiver1@example.com",
		},
	}

	err := SendWithoutAuth(fmt.Sprintf("127.0.0.1:%d", server.PortNumber()), from, to, "Testing subject", "A fantastic email.")

	assert.Nil(t, err)
	assert.Equal(t, 1, len(server.Messages()))
	assert.Equal(t, "MAIL FROM:<john.doe@example.com>", server.Messages()[0].MailfromRequest())
}
