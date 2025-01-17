package text

import (
	"fmt"
	"os"
	"time"

	"gophKeeper/internal/client/model"
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

func (m *Model) Reset() {
	m.Common.Reset()
	m.Data.Reset()
}

func init() {
	model.RegisterModel(&Data{})
}

type Model struct {
	Data *Data `json:"data"`
	model.Common
}

func (m *Model) DataFromFile() (err error) {
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

func (m *Model) Validate(_ ...string) error {
	return model.Validator.Struct(m)
}

func (m *Model) GetPacked() any {
	return m.Data.GetPacked()
}

func (m *Model) GetDst() any {
	return m.Data.GetDst()
}

type Data struct {
	Text string `json:"text" validate:"required" flag:"text,t" default:"" usage:"text field"`
}

func (m *Data) GetPacked() any {
	return m
}

func (m *Data) GetDst() any {
	return m
}

func (m *Data) Reset() {
	m.Text = ""
}
