package data

type Text struct {
	Common
	Text string `json:"text"`
}

func (m *Text) IsValid() (err error) {
	return validate.Struct(m)
}

func (m *Text) Bytes() ([]byte, error) {
	return []byte(m.Text), nil
}
