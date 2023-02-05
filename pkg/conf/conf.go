package conf

import (
	"fmt"

	"github.com/hedge10/airmail/pkg/mail"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SmtpHost	string	`envconfig:"AM_SMTP_HOST" default:"127.0.0.1"`
	SmtpUser	string  `envconfig:"AM_SMTP_USER" default:""`
	SmtpPass	string  `envconfig:"AM_SMTP_PASS" default:""`
	SmtpPort	int		`envconfig:"AM_SMTP_PORT" default:"25"`
	SmtpAuth	string	`envconfig:"AM_SMTP_AUTH_MECHANISM" default:"none"`

	Host		string  `envconfig:"AM_HOST" default:"127.0.0.1"`
	Port		int     `envconfig:"AM_PORT" default:"8081"`

	Env       string  `envconfig:"AM_ENV" default:"prod"`
}

var (
	ErrCannotConnectToSmtp = errors.New("Cannot connect to SMTP server. Please check the connection details.")
)

func isValidAuthMechanism(mechanism string) bool {
    switch mechanism {
	case mail.AUTH_NONE, mail.AUTH_PLAIN, mail.AUTH_LOGIN, mail.AUTH_CRAM_MD5, mail.AUTH_NTLM:
		return true
    }

    return false
}

func (cfg *Config) validate() error {
	if !isValidAuthMechanism(cfg.SmtpAuth) {
		return errors.New(fmt.Sprintf("The authentication mechanism '%s' is unknown. Supported values are: %s, %s, %s, %s, %s.", cfg.SmtpAuth, mail.AUTH_NONE, mail.AUTH_LOGIN, mail.AUTH_PLAIN, mail.AUTH_CRAM_MD5, mail.AUTH_NTLM))
	}

	return nil
}

func (cfg *Config) logging() error {
	if cfg.Env == "dev" {
		log.Debug("Set log level to debug")
		log.SetLevel(log.DebugLevel)
	}

	return nil
}


func New() (*Config, error) {
	config := new(Config)

	err := envconfig.Process("", config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse environment config")
	}

	err = config.validate()
	if err != nil {
		return nil, errors.Wrap(err, "Failed validation of config")
	}

	err = config.logging()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to configure logging")
	}

	log.WithField("env", config.Env).Info("Configuration loaded")

	return config, nil
}
