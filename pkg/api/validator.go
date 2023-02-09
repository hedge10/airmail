package api

import (
	"github.com/go-playground/validator"
)

type(
	CustomValidator struct {
		Validator *validator.Validate
	}
	ValidationError struct {
		Field string
		Message string
	}
)

func MsgForTag(fe validator.FieldError, def error) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return def.Error() // default error
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)

	return err
}
