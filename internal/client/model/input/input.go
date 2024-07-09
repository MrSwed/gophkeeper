package input

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Validate interface {
	Validate() error
}

type Model interface {
	Validate
	Bytes() []byte
	GetKey() string
	GetDescription() *string
	GetFileName() string
}

type Common struct {
	Key         string  `json:"key" validate:"required"`
	Description *string `json:"description"`
	FileName    string  `json:"fileName"`
}

func (c Common) GetKey() string {
	return c.Key
}

func (c Common) GetDescription() *string {
	return c.Description
}

func (c Common) GetFileName() string {
	return c.FileName
}
