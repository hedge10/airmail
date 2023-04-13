package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type (
	Mailgun struct {
		Domain      string
		PrivateKey  string
		UseEuDomain bool
	}
)

func (m Mailgun) Send(e *Email) error {
	mg := mailgun.NewMailgun(m.Domain, m.PrivateKey)
	if m.UseEuDomain {
		mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	}

	// Creating receipients
	to := make([]string, 0, len(e.To))
	for _, t := range e.To {
		to = append(to, string(t.Email))
	}

	message := mg.NewMessage(e.From.Email, e.Subject, e.Message, to...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)

	return err
}
