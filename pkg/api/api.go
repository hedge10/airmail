package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/mail"
)

func sendMail(c echo.Context) error {
	cfg := c.Get("config").(*conf.Config)
	address := fmt.Sprintf("%s:%d", cfg.SmtpHost, cfg.SmtpPort)

	email := new(mail.Email)
	if err := c.Bind(&email); err != nil {
		return c.String(http.StatusBadRequest, "Cannot bind incoming data")
	}

	var sending_err error
	if cfg.SmtpAuth == mail.AUTH_NONE {
		sending_err = mail.SendWithoutAuth(address, email.From, email.To, "", "")
	} else {
		client := mail.CreateClient(cfg.SmtpAuth)
		sending_err = mail.Send(address, client, email.From, email.To, "", "")
	}

	if sending_err == nil {
		return c.JSON(200, "Mail successfully handed over to SMTP server")
	}

	return c.JSON(500, "Mail was not sent")
}
