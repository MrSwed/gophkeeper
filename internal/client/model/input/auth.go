package input

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	Common
	Data AuthData `json:"data"`
}

func (m *Auth) Validate() (err error) {
	return validate.Struct(m)
}

func (m *Auth) Bytes() []byte {
	return []byte(m.Data.Login + ":" + m.Data.Password)
}
