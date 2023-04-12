package mail

import (
	"fmt"
	"strings"

	"github.com/qiniu/qmgo/field"
	log "github.com/sirupsen/logrus"
)

const (
	CT_HTML  string = "text/html"
	CT_PLAIN string = "text/plain"
)

type (
	Meta struct {
		ContentType string
	}
	Party struct {
		Name  string
		Email string
	}
	Email struct {
		field.DefaultField `bson:",inline"`
		Meta               Meta
		From               Party
		To                 []Party
		Cc                 []Party
		Bcc                []Party
		Subject            string
		Message            string
	}
)

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
