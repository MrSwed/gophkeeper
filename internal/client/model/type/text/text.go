package text

import (
	"gophKeeper/internal/client/model"
)

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

type Model struct {
	model.Common
	Data model.Data `json:"data"`
}

func (m *Model) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m.Data.Type(), m.Data)
}

func (m *Model) Type() string {
	return model.GetName(m)
}

func (m *Model) GetData() any {
	return m.Data.GetData()
}

type Data struct {
	Text string `json:"text"`
}

func (m *Data) Type() string {
	return model.GetName(m)
}

func (m *Data) GetData() any {
	return m
}

func init() {
	model.RegisterModel(&Data{})
}
