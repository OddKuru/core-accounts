package vo

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ID struct {
	value string
}

func NewID(id string) (ID, error) {
	result := ID{
		value: id,
	}

	if err := result.Validate(); err != nil {
		return ID{}, errx.WrapWithCode(err, codex.InvalidArgument, "ID.Validate")
	}

	return result, nil
}

func (i ID) Value() string {
	return i.value
}

func (i ID) String() string {
	return i.value
}

func (i ID) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.value,
			validation.Required.Error("value is empty"),
			is.UUIDv4),
	)
}
