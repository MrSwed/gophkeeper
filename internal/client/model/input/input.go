package input

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Model interface {
	Bytes() ([]byte, error)
	Validate() error
}

type Common struct {
	Key         string `json:"key" validate:"required"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	createdDate time.Time
	updatedDate time.Time
}
