package model

import (
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

type Validate interface {
	Validate(fields ...string) error
}
