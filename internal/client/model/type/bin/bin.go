package bin

import (
	"gophKeeper/internal/client/model"
)

type Data struct {
	Bin []byte `json:"bin"`
}

type Model struct {
	model.Common
	Data model.Data `json:"data"`
}

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

func (m *Model) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(model.GetName(m), m.Data)
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
