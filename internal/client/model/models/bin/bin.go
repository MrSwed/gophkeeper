package bin

import (
	"gophKeeper/internal/client/model"
	"reflect"
	"strings"
)

type BinData struct {
	Bin []byte `json:"bin"`
}

type Bin struct {
	model.Common
	Data BinData `json:"data"`
}

var _ model.Model = (*Bin)(nil)

func (m *Bin) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Bin) Bytes() []byte {
	return m.Data.Bin
}

func (m *Bin) Type() string {
	p := strings.Split(reflect.TypeOf(m).PkgPath(), "/")
	return p[len(p)-1]
}
