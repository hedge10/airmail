package mail

import (
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	log "github.com/sirupsen/logrus"
)

const (
	AUTH_NONE     string = "none"
	AUTH_PLAIN    string = "plain"
	AUTH_CRAM_MD5 string = "cram-md5"
	AUTH_LOGIN    string = "login"
	AUTH_NTLM     string = "ntlm"
)

type (
	party struct {
		Name  string `json:"name" form:"name"`
		Email string `json:"email" form:"email" validate:"required,email"`
	}

	Email struct {
		ContentType string  `json:"content-type" form:"content-type"`
		From        party   `json:"from" form:"from" validate:"required,dive"`
		To          []party `json:"to" form:"to" validate:"required,dive,required"`
		Cc          []party `json:"cc" form:"cc" validate:"dive"`
		Bcc         []party `json:"bcc" form:"bcc" validate:"dive"`
		Subject     string  `json:"subject" form:"subject" validate:"required"`
		Message     string  `json:"message" form:"message" validate:"required"`
	}
)

func CreateClient(auth_mechanism string) sasl.Client {
	if auth_mechanism == AUTH_LOGIN {
		log.Debug(fmt.Sprintf("Using %s method for authentication", AUTH_LOGIN))
		return sasl.NewLoginClient("", "")
	}

	log.Debug(fmt.Sprintf("Using %s method for authentication", AUTH_PLAIN))
	return sasl.NewPlainClient("", "", "")
}

func Send(server_address string, client sasl.Client, from party, to []party, subject string, msg string) error {
	sender := buildSender(from)
	receivers := buildReceivers(to)
	message := buildMessage(sender, receivers, subject, msg)

	log.Info("Sending mail to SMTP")
	err := smtp.SendMail(server_address, client, sender, receivers, message)
	if err != nil {
		log.Error(err)
	}

	return err
}

func SendWithoutAuth(server_address string, from party, to []party, subject string, msg string) error {
	c, err := smtp.Dial(server_address)

	fmt.Print(err == nil)

	if err != nil {
		log.Error(err)
		return err
	}
	defer c.Close()

	sender := buildSender(from)
	receivers := buildReceivers(to)
	message := buildMessage(sender, buildRawReceivers(to), subject, msg)

	log.Info("Sending mail to SMTP")
	error := c.SendMail(sender, receivers, message)

	if error != nil {
		log.Error(error)
	}

	return error
}

func buildSender(from party) string {
	sender := from.Email
	if len(from.Name) > 0 {
		sender = fmt.Sprintf("%s<%s>", from.Name, from.Email)
	}

	return sender
}

func buildReceivers(to []party) []string {
	var receivers []string
	for _, el := range to {
		if len(el.Name) > 0 {
			receivers = append(receivers, fmt.Sprintf("%s<%s>", el.Name, el.Email))
		} else {
			receivers = append(receivers, el.Email)
		}
	}

	log.Debug(fmt.Sprintf("Built %d receivers", len(receivers)))

	return receivers
}

func buildRawReceivers(to []party) []string {
	var receivers []string
	for _, el := range to {
		receivers = append(receivers, el.Email)
	}

	log.Debug(fmt.Sprintf("Built %d raw receivers", len(receivers)))

	return receivers
}

func buildMessage(sender string, receipients []string, subject string, message string) *strings.Reader {
	joined_rcpts := strings.Join(receipients, ",")
	from := fmt.Sprintf("From: %s", sender)
	to := fmt.Sprintf("To: %s", joined_rcpts)
	sub := fmt.Sprintf("Subject: %s", subject)

	return strings.NewReader("Content-type: text/html; charset=utf-8" + "\r\n" +
		from + "\r\n" +
		to + "\r\n" +
		sub + "\r\n\r\n" +
		message + "\r\n")
}
