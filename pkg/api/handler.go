package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form"
	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/mail"
	"github.com/hedge10/airmail/pkg/storage"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"github.com/mholt/binding"

	log "github.com/sirupsen/logrus"
)

const (
	CT_FORM = "application/x-www-form-urlencoded"
	CT_JSON = "application/json"
)

type Person struct {
	Name    string `json:"name" form:"name"`
	Address string `json:"address" form:"address"`
}
type MessageRequest struct {
	SenderName    string   `json:"sender-name" form:"sender-name"`
	SenderAddress string   `json:"sender-address" form:"sender-address"`
	To            []Person `json:"to" form:"to"`
	Cc            []Person `json:"cc" form:"cc"`
	Bcc           []Person `json:"bcc" form:"bcc"`
	Subject       string   `json:"subject" form:"subject"`
	Message       string   `json:"message" form:"message"`

	GrecaptchaResponse string `json:"g-recaptcha-response" form:"g-recaptcha-response"`

	ContentType string `json:"_content-type" form:"_content-type"`
	Redirect    string `json:"_redirect" form:"_redirect"`
}

var decoder *form.Decoder

func (mr *MessageRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mr.SenderName: "sender-name",
	}
}

func (mr MessageRequest) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&mr.SenderAddress, validation.Required, is.Email),
		validation.Field(&mr.ContentType, validation.In("html", "text")),
		validation.Field(&mr.Redirect, is.URL),
	}
	err := validation.ValidateStruct(&mr, rules...)
	if err != nil {
		return err
	}

	// Separately validating the receivers, as the built-in logic for nested structs was not that obvious
	// See https://github.com/invopop/validation/issues/3
	for _, el := range mr.To {
		err = validation.Errors{
			"receiver address": validation.Validate(el.Address, validation.Required, is.Email),
		}.Filter()
		if err != nil {
			return err
		}
	}

	return nil
}

func IncomingMessageHandler(config *conf.Config, storage *storage.Storage) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var mr *MessageRequest

		// We can simply check the various content types here, as we cannot get any unknown content type due to the EnforceContentType middleware.
		if r.Header.Get("Content-Type") == CT_FORM {
			if err := r.ParseForm(); err != nil {
				log.Fatal("Cannot parse form values.")
			}
			mr = bindFormData(r.Form)
		}
		if r.Header.Get("Content-Type") == CT_JSON {
			mr = bindJsonData(r)
		}

		// If, for whatever reasons, the binding failed, exit here.
		if mr == nil {
			log.Warn("Binding failed. Skipping further processing.")
			return
		}

		// Validation
		err := mr.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		// Verify Google recaptcha
		if len(config.GrecaptchaSecret) > 0 {
			log.Info("Google Recaptcha active. Trying to validate...")
			c, _ := CreateClient(BaseUri(GOOGLE_SITE_VERIFY))
			r := c.ValidateGrecaptcha(config.GrecaptchaSecret, mr.GrecaptchaResponse, r.RemoteAddr)
			if r != nil {
				log.WithField("grecaptcha_error", r).Debug("Google Recaptcha validation failed.")
				http.Error(w, "Google Recaptcha validation failed.", http.StatusUnprocessableEntity)
				return
			}
		}

		// Start sending email
		from := mail.Party{
			Name:  mr.SenderName,
			Email: mr.SenderAddress,
		}
		address := fmt.Sprintf("%s:%d", config.SmtpHost, config.SmtpPort)
		to := buildReceivers(mr.To)
		cc := buildReceivers(mr.Cc)
		bcc := buildReceivers(mr.Bcc)

		m := mail.Email{
			Connection: mail.Connection{
				Address: address,
			},
			Meta: mail.Meta{
				ContentType: strings.ToLower(mr.ContentType),
			},
			From:    from,
			To:      to,
			Cc:      cc,
			Bcc:     bcc,
			Subject: mr.Subject,
			Message: mr.Message,
		}

		var e error
		if config.SmtpAuth == mail.AUTH_NONE {
			e = m.SendWithoutAuth()
		} else {
			m.Connection.Client = mail.CreateClient(config.SmtpAuth, config.SmtpUser, config.SmtpPass)
			e = m.Send()
		}
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}

		// Successfully sent the message, now store it
		if storage != nil {
			storage.Message = m
			err := storage.Insert()
			if err != nil {
				log.WithField("error", err).Error("Message was not saved.")
			}
		}

		if len(mr.Redirect) > 0 {
			http.Redirect(w, r, mr.Redirect, http.StatusPermanentRedirect)
		}
	}

	return http.HandlerFunc(fn)
}

func bindFormData(values url.Values) (m *MessageRequest) {
	decoder = form.NewDecoder()

	var mr *MessageRequest

	err := decoder.Decode(&mr, values)
	if err != nil {
		log.Fatal(err)
	}

	return mr
}

func bindJsonData(r *http.Request) (m *MessageRequest) {
	mr := new(MessageRequest)

	if errs := binding.Bind(r, mr); errs != nil {
		log.Fatal(errs)
		return
	}

	return mr
}

func buildReceivers(r []Person) []mail.Party {
	var to []mail.Party
	for _, el := range r {
		to = append(to, mail.Party{
			Name:  el.Name,
			Email: el.Address,
		})
	}

	return to
}
