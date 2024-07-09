package input

type Auth struct {
	Common
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m *Auth) Validate() (err error) {
	return validate.Struct(m)
}

func (m *Auth) Bytes() []byte {
	return []byte(m.Login + ":" + m.Password)
}
