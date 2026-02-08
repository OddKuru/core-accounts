package vo

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	result := Email{
		value: value,
	}
	if err := result.Validate(); err != nil {
		return Email{}, errx.WrapWithCode(err, codex.InvalidArgument, "Email.Validate")
	}
	return result, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(
			&e.value,
			validation.Required.Error("value is empty"),
			is.Email.Error("value is not email"),
		),
	)
}
