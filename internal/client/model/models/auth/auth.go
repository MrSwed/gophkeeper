package auth

import (
	"gophKeeper/internal/client/model"
	"reflect"
	"strings"
)

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	model.Common
	Data AuthData `json:"data"`
}

var _ model.Model = (*Auth)(nil)

func (m *Auth) Validate() (err error) {
	return model.Validator.Struct(m)
}

func (m *Auth) Bytes() []byte {
	return []byte(m.Data.Login + ":" + m.Data.Password)
}

func (m *Auth) Type() string {
	p := strings.Split(reflect.TypeOf(m).PkgPath(), "/")
	return p[len(p)-1]
}
