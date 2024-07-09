package input

type TextData struct {
	Text string `json:"text"`
}

type Text struct {
	Common
	Data TextData `json:"data"`
}

func (m *Text) Validate() error {
	return validate.Struct(m)
}

func (m *Text) Bytes() []byte {
	return []byte(m.Data.Text)
}
