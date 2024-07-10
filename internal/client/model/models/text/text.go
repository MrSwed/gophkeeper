package text

import (
	"gophKeeper/internal/client/model"
	"reflect"
	"strings"
)

type TextData struct {
	Text string `json:"text"`
}

type Text struct {
	model.Common
	Data TextData `json:"data"`
}

var _ model.Model = (*Text)(nil)

func (m *Text) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Text) Bytes() []byte {
	return []byte(m.Data.Text)
}
func (m *Text) Type() string {
	p := strings.Split(reflect.TypeOf(m).PkgPath(), "/")
	return p[len(p)-1]
}
