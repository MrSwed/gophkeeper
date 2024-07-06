package input

type ListQuery struct {
	Key         string `json:"key" validate:"omitempty,max=100"`
	Description string `json:"description" validate:"omitempty,max=5000"`
	CreatedAt   string `json:"created_at" validate:"omitempty,datetime"`
	UpdatedAt   string `json:"updated_at" validate:"omitempty,datetime"`
	Limit       uint64 `json:"limit" validate:"omitempty" default:"10"`
	Offset      uint64 `json:"offset" validate:"omitempty"`
}
