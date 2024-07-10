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

func (m *Model) Bytes() []byte {
	return m.Data.Bin
}

func (m *Model) Type() string {
	return model.GetName(m)
}
