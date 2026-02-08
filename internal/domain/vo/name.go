package vo

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginName struct {
	value string
}

func NewLoginName(value string) (LoginName, error) {
	result := LoginName{
		value: value,
	}

	if err := result.Validate(); err != nil {
		return result, errx.WrapWithCode(err, codex.InvalidArgument, "LoginName.Validate")
	}

	return result, nil
}

func (l LoginName) Value() string {
	return l.value
}

func (l LoginName) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.value, validation.Required.Error("value is empty")),
	)
}
