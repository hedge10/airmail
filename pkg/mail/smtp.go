package mail

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/hedge10/airmail/constants"
)

type (
	Smtp struct {
		address string
		auth    sasl.Client
	}
)

func (s Smtp) Send(e *Email) error {
	if s.auth == nil {
		return s.sendWithoutAuth(e)
	}

	return s.send(e)
}

func CreateSmtpAuthClient(auth string, user string, pass string) sasl.Client {
	log.Debug(fmt.Sprintf("Using '%s' method for authentication", auth))

	if auth == constants.AUTH_LOGIN {
		return sasl.NewLoginClient(user, pass)
	}

	return sasl.NewPlainClient("", user, pass)
}

func (s Smtp) sendWithoutAuth(e *Email) error {
	c, err := smtp.Dial(s.address)

	if err != nil {
		log.Error(err)
		return err
	}
	defer c.Close()

	sender := buildSender(e.From)
	receivers := buildReceivers(e.To)
	message := buildMessage(sender, buildRawReceivers(e.To), buildRawReceivers(e.Cc), buildRawReceivers(e.Bcc), e.Subject, e.Message, e.Meta.ContentType)

	log.Info("Sending mail to SMTP123")
	log.Info(s.address)
	error := c.SendMail(sender, receivers, message)

	if err != nil {
		log.WithField("smtp_error", err).Error("Sending mail failed")
	}

	return error
}

func (s Smtp) send(e *Email) error {
	sender := buildSender(e.From)
	receivers := buildReceivers(e.To)
	message := buildMessage(sender, receivers, buildRawReceivers(e.Cc), buildRawReceivers(e.Bcc), e.Subject, e.Message, e.Meta.ContentType)

	if s.auth == nil {
		log.Warn("no authentication client initialized.")
		return errors.New("no authentication client initialized")
	}

	err := smtp.SendMail(s.address, s.auth, sender, receivers, message)
	if err != nil {
		log.WithField("smtp_error", err).Error("Sending mail failed")
	}

	return err
}
