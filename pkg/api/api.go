package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/mail"
)

type (
	AirmailResponse struct {
		Code int16
		Message string
		Errors []ValidationError
	}
)

func sendMail(c echo.Context) error {
	cfg := c.Get("config").(*conf.Config)
	address := fmt.Sprintf("%s:%d", cfg.SmtpHost, cfg.SmtpPort)

	email := new(mail.Email)
	if err := c.Bind(&email); err != nil {
		return c.String(http.StatusBadRequest, "Cannot bind incoming data")
	}

	// Validate incoming data
	if err := c.Validate(email); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ValidationError, len(ve))
			for i, fe := range ve {
				out[i] = ValidationError{fe.Field(), MsgForTag(fe, err)}
			}

			r := AirmailResponse {
				Code: http.StatusBadRequest,
				Message: "An error occured",
				Errors: out,
			}

			return c.JSONPretty(http.StatusBadRequest, r, " ")
		}
	}

	var sending_err error
	if cfg.SmtpAuth == mail.AUTH_NONE {
		sending_err = mail.SendWithoutAuth(address, email.From, email.To, email.Subject, email.Message)
	} else {
		client := mail.CreateClient(cfg.SmtpAuth)
		sending_err = mail.Send(address, client, email.From, email.To, email.Subject, email.Message)
	}

	if sending_err == nil {
		return c.JSON(http.StatusOK, AirmailResponse {
			Code: http.StatusOK,
			Message: "Mail successfully handed over to SMTP server.",
		})
	}

	return c.JSON(http.StatusInternalServerError, AirmailResponse {
		Code: http.StatusInternalServerError,
		Message: "An internal error occured. Mail was not sent.",
	})
}
