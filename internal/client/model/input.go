package model

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

type Validate interface {
	Validate() error
}

type Packed struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewPackedBytes(t string, d any) ([]byte, error) {
	return json.Marshal(Packed{
		Type: t,
		Data: d,
	})
}

type ListQuery struct {
	Key         string `json:"key" validate:"omitempty,max=100"`
	Description string `json:"description" validate:"omitempty,max=5000"`
	CreatedAt   string `json:"created_at" validate:"omitempty,datetime"`
	UpdatedAt   string `json:"updated_at" validate:"omitempty,datetime"`
	Limit       uint64 `json:"limit" validate:"omitempty" default:"10"`
	Offset      uint64 `json:"offset" validate:"omitempty"`
}

func (m *ListQuery) Validate() (err error) {
	return Validator.Struct(m)
}
