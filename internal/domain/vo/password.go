package vo

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	MinPasswordLen = 8
	MaxPasswordLen = 64
)

type (
	Password struct {
		value string
	}

	HashedPassword struct {
		value []byte
	}
)

func NewPassword(value string) (Password, error) {
	result := Password{
		value: value,
	}

	if err := result.Validate(); err != nil {
		return Password{}, errx.WrapWithCode(err, codex.InvalidArgument, "Password.Validate")
	}

	return result, nil
}

func NewHashedPassword(value []byte) HashedPassword {
	return HashedPassword{
		value: value,
	}
}

func (p Password) Value() string {
	return p.value
}

func (p Password) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.value, validation.Required.Error("value is empty")),
		validation.Field(&p.value, validation.Length(MinPasswordLen, MaxPasswordLen)),
	)
}

func (p HashedPassword) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.value, validation.Required.Error("value is empty")),
	)
}

func (p HashedPassword) Value() []byte {
	return p.value
}
