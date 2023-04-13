package conf

import (
	"fmt"

	"github.com/hedge10/airmail/constants"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	MailService string `envconfig:"AM_MAIL_SERVICE" default:"smtp"`

	MailgunDomain      string `envconfig:"AM_MAILGUN_DOMAIN"`
	MailgunKey         string `envconfig:"AM_MAILGUN_PRIVATE_KEY"`
	MailgunUseEuDomain bool   `envconfig:"AM_MAILGUN_USE_EU_DOMAIN" default:"false"`

	SmtpHost string `envconfig:"AM_SMTP_HOST" default:"127.0.0.1"`
	SmtpUser string `envconfig:"AM_SMTP_USER" default:""`
	SmtpPass string `envconfig:"AM_SMTP_PASS" default:""`
	SmtpPort int    `envconfig:"AM_SMTP_PORT" default:"25"`
	SmtpAuth string `envconfig:"AM_SMTP_AUTH_MECHANISM" default:"none"`

	Host string `envconfig:"AM_HOST" default:""`
	Port int    `envconfig:"AM_PORT" default:"9900"`

	Debug bool   `envconfig:"AM_DEBUG" default:"false"`
	Env   string `envconfig:"AM_ENV" default:"dev"`

	GrecaptchaSecret string `envconfig:"AM_GRECAPTCHA_SECRET" default:""`

	UseStorage        bool   `envconfig:"AM_USE_STORAGE" default:"false"`
	StorageType       string `envconfig:"AM_STORAGE_TYPE" default:"mongodb"`
	MongoDbHost       string `envconfig:"AM_MONGODB_HOST" default:"localhost"`
	MongoDbPort       int    `envconfig:"AM_MONGODB_PORT" default:"27017"`
	MongoDbDatabase   string `envconfig:"AM_MONGODB_DB" default:"airmail"`
	MongoDbCollection string `envconfig:"AM_MONGODB_COLLECTION" default:"messages"`
	MongoDbUsername   string `envconfig:"AM_MONGODB_USERNAME" default:""`
	MongoDbPassword   string `envconfig:"AM_MONGODB_PASSWORD" default:""`
}

var (
	ErrCannotConnectToSmtp = errors.New("Cannot connect to SMTP server. Please check the connection details.")
)

func isValidAuthMechanism(cfg *Config) error {
	switch cfg.SmtpAuth {
	case constants.AUTH_NONE, constants.AUTH_PLAIN, constants.AUTH_LOGIN:
		return nil
	}

	return errors.New(fmt.Sprintf("The authentication mechanism '%s' is unknown. Supported values are: %s, %s, %s.", cfg.SmtpAuth, constants.AUTH_NONE, constants.AUTH_PLAIN, constants.AUTH_LOGIN))
}

func isValidMailService(cfg *Config) error {
	switch cfg.MailService {
	case constants.MAIL_SERVICE_SMTP, constants.MAIL_SERVICE_MAILGUN:
		return nil
	}

	return errors.New(fmt.Sprintf("The mail service '%s' is unknown. Supported values are: %s, %s.", cfg.MailService, constants.MAIL_SERVICE_MAILGUN, constants.MAIL_SERVICE_SMTP))
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
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info(fmt.Sprintf("Log level set to %s", log.GetLevel()))

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
