package input

type Bin struct {
	Common
	Bin []byte `json:"bin"`
}

func (m *Bin) Validate() error {
	return validate.Struct(m)
}

func (m *Bin) Bytes() []byte {
	return m.Bin
}
