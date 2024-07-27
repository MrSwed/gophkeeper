package auth

import (
	"gophKeeper/internal/client/model"
)

type Data struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Model struct {
	model.Common
	Data *Data `json:"data"`
}

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

func (m *Model) Validate(fields ...string) error {
	if len(fields) == 0 {
		return model.Validator.Struct(m)
	} else {
		return model.Validator.StructPartial(m, fields...)
	}
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m)
}

func (m *Model) GetData() any {
	return m.Data.GetData()
}

func (m *Data) GetData() any {
	return m
}

func init() {
	model.RegisterModel(&Data{})
}
