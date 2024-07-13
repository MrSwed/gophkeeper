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

func (m *Model) Validate() (err error) {
	return model.Validator.Struct(m)
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
