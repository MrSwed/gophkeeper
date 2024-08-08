package text

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"os"
	"time"
)

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

func New() *Model {
	return &Model{
		Common: model.Common{},
		Data:   &Data{},
	}
}

func init() {
	model.RegisterModel(&Data{})
}

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

func (m *Model) GetKey() string {
	if m.Key == "" {
		m.Key = fmt.Sprintf("%s-%s", model.GetName(m), time.Now().Format("2006-01-02-15-04-05"))
	}
	return m.Key
}

func (m *Model) Validate(fields ...string) error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m)
}

func (m *Model) GetPacked() any {
	return m.Data.GetPacked()
}

func (m *Model) GetDst() any {
	return m.Data.GetDst()
}

type Data struct {
	Text string `json:"text"`
}

func (m *Data) GetPacked() any {
	return m
}

func (m *Data) GetDst() any {
	return m
}
