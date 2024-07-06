package data

type Auth struct {
	Common
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m *Auth) IsValid() (err error) {
	return validate.Struct(m)
}

func (m *Auth) Bytes() ([]byte, error) {
	return []byte(m.Login + ":" + m.Password), nil
}
