package model

import "github.com/google/uuid"

type ListQuery struct {
	UserID  uuid.UUID
	Limit   uint64 `json:"limit" validate:"omitempty" default:"10"`
	Offset  uint64 `json:"offset" validate:"omitempty"`
	OrderBy string `json:"orderBy" validate:"omitempty,oneof=key created_at updated_at 'key desc' 'created_at desc' 'updated_at desc'"`
}

func (m *ListQuery) Validate() error {
	if m == nil {
		return nil
	}
	return Validator.Struct(m)
}

// todo:  custom validator for dates range
