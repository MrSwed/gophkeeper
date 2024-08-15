package auth

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"time"
)

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

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

func init() {
	model.RegisterModel(&Data{})
}

func (m *Model) Reset() {
	m.Common.Reset()
	(*m).Data.Reset()
}

func (m *Model) GetKey() string {
	if m.Key == "" {
		m.Key = fmt.Sprintf("%s-%s", model.GetName(m), time.Now().Format("2006-01-02-15-04-05"))
	}
	return m.Key
}

func (m *Model) Validate(fields ...string) error {
	if len(fields) == 0 {
		return model.Validator.Struct(m)
	} else {
		return model.Validator.StructPartial(m, fields...)
	}
}

func (m *Model) GetPacked() any {
	return m.Data.GetPacked()
}

func (m *Model) GetDst() any {
	return m.Data.GetDst()
}

type Data struct {
	Login    string `json:"login" flag:"login,l" default:"" usage:"login field"`
	Password string `json:"password" flag:"password,p" default:"" usage:"password field"`
}

func (m *Data) GetPacked() any {
	return m
}

func (m *Data) GetDst() any {
	return m
}

func (m *Data) Reset() {
	(*m).Password = ""
	(*m).Login = ""
}
