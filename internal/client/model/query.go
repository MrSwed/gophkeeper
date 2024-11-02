package model

type ListQuery struct {
	Key         string `json:"key" validate:"omitempty,max=100" flag:"key,k" usage:"search by key"`
	Description string `json:"description" validate:"omitempty,max=5000" flag:"description,d" usage:"search by description"`
	CreatedAt   string `json:"created_at" validate:"omitempty,datetime=2006-01-02 15:04:05" flag:"created,c" usage:"search by created_at"`
	UpdatedAt   string `json:"updated_at" validate:"omitempty,datetime=2006-01-02 15:04:05" flag:"updated,u" usage:"search by updated_at"`
	SyncAt      string `json:"sync_at" validate:"omitempty,datetime=2006-01-02 15:04:05" flag:"sync,s" usage:"get all earlier than the sync_at"`
	Limit       uint64 `json:"limit" validate:"omitempty" default:"10" flag:"limit,l" usage:"set limit"`
	Offset      uint64 `json:"offset" validate:"omitempty" flag:"offset,o" usage:"set offset"`
	OrderBy     string `json:"orderBy" validate:"omitempty,oneof=key created_at updated_at sync_at 'key desc' 'created_at desc' 'updated_at desc' 'sync_at desc'" flag:"order-by,b" usage:"set order by"`
	Deleted     bool   `json:"deleted" flag:"deleted" usage:"show deleted"`
}

func (m *ListQuery) Validate() (err error) {
	return Validator.Struct(m)
}
