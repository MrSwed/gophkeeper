package model

type ListQuery struct {
	Key         string `json:"key" validate:"omitempty,max=100" flag:"key" usage:"search by key"`
	Description string `json:"description" validate:"omitempty,max=5000" flag:"description" usage:"search by description"`
	CreatedAt   string `json:"created_at" validate:"omitempty,datetime" flag:"created" usage:"search by created_at"`
	UpdatedAt   string `json:"updated_at" validate:"omitempty,datetime" flag:"updated" usage:"search by updated_at"`
	Limit       uint64 `json:"limit" validate:"omitempty" default:"10" flag:"limit" usage:"set limit"`
	Offset      uint64 `json:"offset" validate:"omitempty" flag:"offset" usage:"set offset"`
}

func (m *ListQuery) Validate() (err error) {
	return Validator.Struct(m)
}
