package conf

import (
	"fmt"
	"testing"

	"github.com/hedge10/airmail/constants"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValidateMailService(t *testing.T) {
	type test struct {
		input string
		want  error
	}

	tests := []test{
		{input: "smtp", want: nil},
		{input: "mailgun", want: nil},
		{input: "ses", want: nil},
		{input: "unknown", want: fmt.Errorf("The mail service '%s' is unknown. Supported values are: %s, %s, %s.", "unknown", constants.MAIL_SERVICE_MAILGUN, constants.MAIL_SERVICE_SMTP, constants.MAIL_SERVICE_SES)},
	}

	for _, tc := range tests {
		c := &Config{
			MailService: tc.input,
		}

		got := isValidMailService(c)
		if tc.want != nil {
			assert.Equal(t, tc.want.Error(), got.Error())
		} else {
			assert.Nil(t, got)
		}
	}
}

func TestValidateAuthMechanism(t *testing.T) {
	type test struct {
		input string
		want  error
	}

	tests := []test{
		{input: "none", want: nil},
		{input: "plain", want: nil},
		{input: "login", want: nil},
		{input: "unknown", want: fmt.Errorf("The authentication mechanism '%s' is unknown. Supported values are: %s, %s, %s.", "unknown", constants.AUTH_NONE, constants.AUTH_PLAIN, constants.AUTH_LOGIN)},
	}

	for _, tc := range tests {
		c := &Config{
			SmtpAuth: tc.input,
		}

		got := isValidAuthMechanism(c)
		if tc.want != nil {
			assert.Equal(t, tc.want.Error(), got.Error())
		} else {
			assert.Nil(t, got)
		}
	}
}

func TestValidateLogLevel(t *testing.T) {
	type test struct {
		input bool
		level string
	}

	tests := []test{
		{input: true, level: "debug"},
		{input: false, level: "info"},
	}

	for _, tc := range tests {
		c := &Config{
			Debug: tc.input,
		}

		err := c.logging()

		assert.Nil(t, err)
		assert.Equal(t, tc.level, log.GetLevel().String())
	}
}
