package bin

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

func init() {
	model.RegisterModel(&Data{})
}

type Model struct {
	model.Common
	Data *Data `json:"data"`
}

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

func (m *Model) DataFromFile() (err error) {
	if m.FileName != "" {
		if m.Data == nil {
			m.Data = &Data{}
		}
		m.Data.Bin, err = os.ReadFile(m.FileName)
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
	Bin []byte `json:"bin" validate:"required"`
}

func (m *Data) GetPacked() any {
	return m
}

func (m *Data) GetDst() any {
	return m
}

func (m *Data) Reset() {
	(*m).Bin = nil
}
