package model

import (
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

type Validate interface {
	Validate(fields ...string) error
}

func ValidateStruct[T Validate](p T, fields ...string) error {
	if len(fields) == 0 {
		return Validator.Struct(p)
	} else {
		return Validator.StructPartial(p, fields...)
	}
}
