package mail

import (
	"errors"
	"fmt"

	"github.com/emersion/go-sasl"
	"github.com/hedge10/airmail/constants"
	"github.com/hedge10/airmail/pkg/conf"
)

type Transfer interface {
	Send(e *Email) error
}

func CreateTransfer(config *conf.Config) (Transfer, error) {
	if config.MailService == constants.MAIL_SERVICE_MAILGUN {
		return &Mailgun{
			Domain:      config.MailgunDomain,
			PrivateKey:  config.MailgunKey,
			UseEuDomain: config.MailgunUseEuDomain,
		}, nil
	}

	if config.MailService == constants.MAIL_SERVICE_SMTP {
		var client sasl.Client
		if (config.SmtpAuth == constants.AUTH_PLAIN) || (config.SmtpAuth == constants.AUTH_LOGIN) {
			client = CreateSmtpAuthClient(config.SmtpAuth, config.SmtpUser, config.SmtpPass)
		}

		return &Smtp{
			address: fmt.Sprintf("%s:%d", config.SmtpHost, config.SmtpPort),
			auth:    client,
		}, nil

	}

	if config.MailService == constants.MAIL_SERVICE_SES {
		return &Ses{
			AwsAccessKeyId:     config.SesAwsAccessKeyId,
			AwsSecretAccessKey: config.SesAwsSecretAccessKey,
			Region:             config.SesRegion,
		}, nil
	}

	return nil, errors.New("cannot create transfer")
}
