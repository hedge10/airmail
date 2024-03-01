package conf

import (
	"context"
	"fmt"

	"github.com/hedge10/airmail/constants"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	MailService string `env:"AM_MAIL_SERVICE,default=smtp"`

	SesRegion             string `env:"AM_SES_REGION,default=$AWS_REGION"`
	SesAwsAccessKeyId     string `env:"AM_SES_AWS_ACCESS_KEY_ID,default=$AWS_ACCESS_KEY_ID"`
	SesAwsSecretAccessKey string `env:"AM_SES_AWS_SECRET_ACCESS_KEY,default=$AWS_SECRET_ACCESS_KEY"`

	MailgunDomain      string `env:"AM_MAILGUN_DOMAIN"`
	MailgunKey         string `env:"AM_MAILGUN_PRIVATE_KEY"`
	MailgunUseEuDomain bool   `env:"AM_MAILGUN_USE_EU_DOMAIN,default=false"`

	SmtpHost string `env:"AM_SMTP_HOST,default=127.0.0.1"`
	SmtpUser string `env:"AM_SMTP_USER"`
	SmtpPass string `env:"AM_SMTP_PASS"`
	SmtpPort int    `env:"AM_SMTP_PORT,default=25"`
	SmtpAuth string `env:"AM_SMTP_AUTH_MECHANISM,default=none"`

	AuthDisabled       bool   `env:"AM_AUTH_DISABLED,default=false"`
	AuthToken          string `env:"AM_AUTH_TOKEN"`
	CorsAllowOrigin    string `env:"AM_CORS_ALLOW_ORIGIN"`
	CorsAllowedHeaders string `env:"AM_CORS_ALLOWED_HEADERS"`

	Host string `env:"AM_HOST"`
	Port int    `env:"AM_PORT,default=9900"`

	Debug bool   `env:"AM_DEBUG,default=false"`
	Env   string `env:"AM_ENV,default=dev"`
}

func isValidAuthMechanism(cfg *Config) error {
	switch cfg.SmtpAuth {
	case constants.AUTH_NONE, constants.AUTH_PLAIN, constants.AUTH_LOGIN:
		return nil
	}

	return errors.New(fmt.Sprintf("The authentication mechanism '%s' is unknown. Supported values are: %s, %s, %s.", cfg.SmtpAuth, constants.AUTH_NONE, constants.AUTH_PLAIN, constants.AUTH_LOGIN))
}

func isValidMailService(cfg *Config) error {
	switch cfg.MailService {
	case constants.MAIL_SERVICE_SMTP, constants.MAIL_SERVICE_MAILGUN, constants.MAIL_SERVICE_SES:
		return nil
	}

	return errors.New(fmt.Sprintf("The mail service '%s' is unknown. Supported values are: %s, %s, %s.", cfg.MailService, constants.MAIL_SERVICE_MAILGUN, constants.MAIL_SERVICE_SMTP, constants.MAIL_SERVICE_SES))
}

func (cfg *Config) validate() error {
	err := isValidAuthMechanism(cfg)
	if err != nil {
		return err
	}

	err = isValidMailService(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) logging() error {
	log.SetLevel(log.InfoLevel)

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info(fmt.Sprintf("Log level set to %s", log.GetLevel()))

	return nil
}

func New() (*Config, error) {
	config := new(Config)

	err := envconfig.Process(context.Background(), config)
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
