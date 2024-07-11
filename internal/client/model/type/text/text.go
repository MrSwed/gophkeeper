package text

import (
	"gophKeeper/internal/client/model"
)

type ModelData struct {
	Text string `json:"text"`
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
	return model.NewPackedBytes(m.Type(), m.Data)
}

func (m *Model) Type() string {
	return model.GetName(m)
}
