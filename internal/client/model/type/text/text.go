package text

import (
	"gophKeeper/internal/client/model"
	"os"
)

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

type Model struct {
	model.Common
	Data *Data `json:"data"`
}

func (m *Model) GetFile() (err error) {
	if m.FileName != "" {
		if m.Data == nil {
			m.Data = &Data{}
		}

		var b []byte
		b, err = os.ReadFile(m.FileName)
		m.Data.Text = string(b)
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

type Data struct {
	Text string `json:"text"`
}

func (m *Data) GetData() any {
	return m
}

func init() {
	model.RegisterModel(&Data{})
}
