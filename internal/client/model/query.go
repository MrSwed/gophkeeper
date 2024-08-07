package model

type ListQuery struct {
	Key         string `json:"key" validate:"omitempty,max=100" flag:"key,k" usage:"search by key"`
	Description string `json:"description" validate:"omitempty,max=5000" flag:"description,d" usage:"search by description"`
	CreatedAt   string `json:"created_at" validate:"omitempty,datetime" flag:"created,c" usage:"search by created_at"`
	UpdatedAt   string `json:"updated_at" validate:"omitempty,datetime" flag:"updated,u" usage:"search by updated_at"`
	Limit       uint64 `json:"limit" validate:"omitempty" default:"10" flag:"limit,l" usage:"set limit"`
	Offset      uint64 `json:"offset" validate:"omitempty" flag:"offset,o" usage:"set offset"`
}

func (m *ListQuery) Validate() (err error) {
	return Validator.Struct(m)
}
