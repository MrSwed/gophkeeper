package auth

import (
	"gophKeeper/internal/client/model"
)

type ModelData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Model struct {
	model.Common
	Data ModelData `json:"data"`
}

var _ model.Model = (*Model)(nil)

func (m *Model) Validate() (err error) {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m.Type(), m.Data)
}

func (m *Model) Type() string {
	return model.GetName(m)
}
