package input

type BinData struct {
	Bin []byte `json:"bin"`
}

type Bin struct {
	Common
	Data BinData `json:"data"`
}

func (m *Bin) Validate() error {
	return validate.Struct(m)
}

func (m *Bin) Bytes() []byte {
	return m.Data.Bin
}
