package mail

import (
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	log "github.com/sirupsen/logrus"
)

const (
	AUTH_NONE string = "none"
	AUTH_PLAIN string = "plain"
	AUTH_CRAM_MD5 string = "cram-md5"
	AUTH_LOGIN string = "login"
	AUTH_NTLM string = "ntlm"
)

type(
	party struct {
		Name string `json:"name"`
		Email	string `json:"email" validate:"required,email"`
	}

	Email struct {
		From party `json:"from" validate:"required,dive"`
		To []party `json:"to" validate:"required,dive,required"`
		Cc []party `json:"cc" validate:"dive"`
		Bcc []party `json:"bcc" validate:"dive"`
		Subject string `json:"subject" validate:"required"`
		Message string `json:"message" validate:"required"`
	}
)


func CreateClient(auth_mechanism string) sasl.Client {
	if auth_mechanism == AUTH_LOGIN {
		log.Info(fmt.Sprintf("Using %s method for authentication", AUTH_LOGIN))
		return sasl.NewLoginClient("", "")
	}

	log.Info(fmt.Sprintf("Using %s method for authentication", AUTH_PLAIN))
	return sasl.NewPlainClient("", "", "")
}

func Send(server_address string, client sasl.Client, from party, to []party, subject string, msg string) error {
	sender := fmt.Sprintf("%s<%s>", from.Name, from.Email)
	receivers := buildReceivers(to)
	message := buildMessage(receivers, subject, msg)

	log.Info("Sending mail to SMTP")
	err := smtp.SendMail(server_address, client, sender, receivers, message)
	if err != nil {
		log.Error(err)
	}

	return err
}

func SendWithoutAuth(server_address string, from party, to []party, subject string, msg string) error {
	c, err := smtp.Dial(server_address)
	if err != nil {
		return err
	}
	defer c.Close()

	sender := fmt.Sprintf("%s<%s>", from.Name, from.Email)
	receivers := buildReceivers(to)
	message := buildMessage(buildRawReceivers(to), subject, msg)



	log.Info("Sending mail to SMTP")
	error := c.SendMail(sender, receivers, message)
	if error != nil {
		log.Error(error)
	}

	return error
}

func buildReceivers(to []party) []string {
	var receivers []string
	for _, el := range to {
		receivers = append(receivers, fmt.Sprintf("%s<%s>", el.Name, el.Email))
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

func buildMessage(receipients []string, subject string, message string) *strings.Reader {
	joined_rcpts := strings.Join(receipients, ",")
	to := fmt.Sprintf("To: %s\r\n", joined_rcpts)
	sub := fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	msg := fmt.Sprintf("%s\r\n", message)

	log.Info("MESSAGE--")

	log.Info("--MESSAGE")

	return strings.NewReader(fmt.Sprintf(to + sub + msg))
}
