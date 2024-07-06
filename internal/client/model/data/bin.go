package data

type Bin struct {
	Common
	Bin []byte `json:"bin"`
}

func (m *Bin) IsValid() (err error) {
	return validate.Struct(m)
}

func (m *Bin) Bytes() ([]byte, error) {
	return m.Bin, nil
}
