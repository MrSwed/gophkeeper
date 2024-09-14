package model

type ListQuery struct {
	Limit  uint64 `json:"limit" validate:"omitempty" default:"10"`
	Offset uint64 `json:"offset" validate:"omitempty"`
}

func (m *ListQuery) Validate() error {
	if m == nil {
		return nil
	}
	return Validator.Struct(m)
}

// todo:  custom validator for dates range
