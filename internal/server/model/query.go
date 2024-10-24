package model

import "github.com/google/uuid"

type ListQuery struct {
	UserID  uuid.UUID
	OrderBy string `json:"orderBy" validate:"omitempty,oneof=key created_at updated_at 'key desc' 'created_at desc' 'updated_at desc'"`
	Offset  uint64 `json:"offset" validate:"omitempty"`
	Limit   uint64 `json:"limit" validate:"omitempty" default:"10"`
}

func (m *ListQuery) Validate() error {
	if m == nil {
		return nil
	}
	return Validator.Struct(m)
}
