package model

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

type Validate interface {
	Validate() error
}

type Model interface {
	Validate
	Bytes() (b []byte, err error)
	GetKey() string
	GetDescription() *string
	GetFileName() string
	Type() string
}

type Common struct {
	Key         string  `json:"key" validate:"required"`
	Description *string `json:"description"`
	FileName    string  `json:"fileName"`
}

type Packed struct {
	Type string `json:"type"`
	Data any    `json:"data"`
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

func NewPackedBytes(t string, d any) ([]byte, error) {
	return json.Marshal(Packed{
		Type: t,
		Data: d,
	})
}
