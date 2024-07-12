package bin

import (
	"gophKeeper/internal/client/model"
)

type ModelData struct {
	Bin []byte `json:"bin"`
}

type Model struct {
	model.Common
	Data ModelData `json:"data"`
}

var _ model.Model = (*Model)(nil)

func (m *Model) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m.Data.Type(), m.Data)
}

func (m *ModelData) Type() string {
	return model.GetName(m)
}

func (m *ModelData) GetData() any {
	return m
}

func init() {
	model.RegisterModel(&ModelData{})
}
