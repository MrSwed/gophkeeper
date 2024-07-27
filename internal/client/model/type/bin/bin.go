package bin

import (
	"gophKeeper/internal/client/model"
	"os"
)

type Data struct {
	Bin []byte `json:"bin"`
}

type Model struct {
	model.Common
	Data *Data `json:"data"`
}

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

func (m *Model) GetFile() (err error) {
	if m.FileName != "" {
		if m.Data == nil {
			m.Data = &Data{}
		}
		m.Data.Bin, err = os.ReadFile(m.FileName)
		return
	}
	return
}

func (m *Model) Validate(fields ...string) error {
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
