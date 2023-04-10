package mail

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/qiniu/qmgo/field"
	log "github.com/sirupsen/logrus"
)

const (
	AUTH_NONE  string = "none"
	AUTH_PLAIN string = "plain"
	AUTH_LOGIN string = "login"

	CT_HTML  string = "text/html"
	CT_PLAIN string = "text/plain"
)

type (
	Connection struct {
		Address string
		Client  sasl.Client
	}

	Meta struct {
		ContentType string
	}

	Party struct {
		Name  string
		Email string
	}

	Email struct {
		field.DefaultField `bson:",inline"`
		Connection         Connection
		Meta               Meta
		From               Party
		To                 []Party
		Cc                 []Party
		Bcc                []Party
		Subject            string
		Message            string
	}
)

func CreateClient(auth string, user string, pass string) sasl.Client {
	log.Debug(fmt.Sprintf("Using '%s' method for authentication", auth))

	if auth == AUTH_LOGIN {
		return sasl.NewLoginClient(user, pass)
	}

	return sasl.NewPlainClient("", user, pass)
}

func (e Email) Send() error {
	sender := buildSender(e.From)
	receivers := buildReceivers(e.To)
	message := buildMessage(sender, receivers, buildRawReceivers(e.Cc), buildRawReceivers(e.Bcc), e.Subject, e.Message, e.Meta.ContentType)

	if e.Connection.Client == nil {
		log.Warn("no authentication client initialized.")
		return errors.New("no authentication client initialized")
	}

	err := smtp.SendMail(e.Connection.Address, e.Connection.Client, sender, receivers, message)
	if err != nil {
		log.WithField("smtp_error", err).Error("Sending mail failed")
	}

	return err
}

func (e Email) SendWithoutAuth() error {
	c, err := smtp.Dial(e.Connection.Address)

	if err != nil {
		log.Error(err)
		return err
	}
	defer c.Close()

	sender := buildSender(e.From)
	receivers := buildReceivers(e.To)
	message := buildMessage(sender, buildRawReceivers(e.To), buildRawReceivers(e.Cc), buildRawReceivers(e.Bcc), e.Subject, e.Message, e.Meta.ContentType)

	log.Info("Sending mail to SMTP123")
	log.Info(e.Connection.Address)
	error := c.SendMail(sender, receivers, message)

	if err != nil {
		log.WithField("smtp_error", err).Error("Sending mail failed")
	}

	return error
}

func buildSender(from Party) string {
	sender := from.Email
	if len(from.Name) > 0 {
		sender = fmt.Sprintf("%s<%s>", from.Name, from.Email)
	}

	return sender
}

func buildReceivers(to []Party) []string {
	var receivers []string
	for _, el := range to {
		if len(el.Name) > 0 {
			receivers = append(receivers, fmt.Sprintf("%s<%s>", el.Name, el.Email))
		} else {
			receivers = append(receivers, el.Email)
		}
	}

	log.WithField("receivers", receivers).Debug(fmt.Sprintf("Built %d receivers", len(receivers)))

	return receivers
}

func buildRawReceivers(to []Party) []string {
	var receivers []string
	for _, el := range to {
		receivers = append(receivers, el.Email)
	}

	log.WithField("raw_receivers", receivers).Debug(fmt.Sprintf("Built %d raw receivers", len(receivers)))

	return receivers
}

func buildMessage(sender string, receipients []string, cc_receipients []string, bcc_receipients []string, subject string, message string, ct string) *strings.Reader {
	var ctype string
	if ct == "html" {
		ctype = fmt.Sprintf("Content-type: %s;", CT_HTML)
	} else {
		ctype = fmt.Sprintf("Content-type: %s;", CT_PLAIN)
	}

	to := fmt.Sprintf("To: %s", strings.Join(receipients, ","))
	cc := fmt.Sprintf("Cc: %s", strings.Join(cc_receipients, ","))
	bcc := fmt.Sprintf("Bcc: %s", strings.Join(bcc_receipients, ","))

	from := fmt.Sprintf("From: %s", sender)
	sub := fmt.Sprintf("Subject: %s", subject)

	return strings.NewReader(ctype + "charset=utf-8" + "\r\n" +
		from + "\r\n" +
		to + "\r\n" +
		cc + "\r\n" +
		bcc + "\r\n" +
		sub + "\r\n\r\n" +
		message + "\r\n")
}
