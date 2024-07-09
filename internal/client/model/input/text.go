package input

type Text struct {
	Common
	Text string `json:"text"`
}

func (m *Text) Validate() error {
	return validate.Struct(m)
}

func (m *Text) Bytes() []byte {
	return []byte(m.Text)
}
