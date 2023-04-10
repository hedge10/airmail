package conf

import (
	"fmt"
	"os"
	"testing"

	"github.com/hedge10/airmail/pkg/mail"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValidateAuthMechanism(t *testing.T) {
	os.Setenv("AM_SMTP_AUTH_MECHANISM", "unknown-auth")
	defer os.Unsetenv("AM_SMTP_AUTH_MECHANISM")

	conf, err := New()

	assert.Nil(t, conf)
	assert.NotNil(t, err)
	assert.Equal(
		t,
		fmt.Sprintf("Failed validation of config: The authentication mechanism '%s' is unknown. Supported values are: %s, %s, %s.", "unknown-auth", mail.AUTH_NONE, mail.AUTH_LOGIN, mail.AUTH_PLAIN),
		err.Error(),
	)
}

func TestDefaultInfoLogLevel(t *testing.T) {
	os.Setenv("AM_DEBUG", "false")
	defer os.Unsetenv("AM_DEBUG")

	conf, err := New()

	assert.NotNil(t, conf)
	assert.Nil(t, err)
	assert.Equal(t, "info", log.GetLevel().String())
}

func TestSetDebugLogLevel(t *testing.T) {
	os.Setenv("AM_DEBUG", "true")
	defer os.Unsetenv("AM_DEBUG")

	conf, err := New()

	assert.NotNil(t, conf)
	assert.Nil(t, err)
	assert.Equal(t, "debug", log.GetLevel().String())
}
