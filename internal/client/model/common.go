package model

import (
	"encoding/json"
	"path/filepath"
)

type Common struct {
	Key         string  `json:"key" validate:"required"`
	Description *string `json:"description"`
	FileName    string  `json:"fileName"`
}

func (c Common) GetKey() string {
	return c.Key
}

func (c Common) GetDescription() *string {
	return c.Description
}
func (c Common) GetFileName() string {
	return c.FileName
}

type Packed struct {
	Type     string `json:"type"`
	Data     any    `json:"data"`
	FileName string `json:"fileName,omitempty"`
}

func NewPackedBytes(m Model) ([]byte, error) {
	if m, ok := m.(GetFile); ok {
		err := m.GetFile()
		if err != nil {
			return nil, err
		}
	}
	p := Packed{
		Type: GetName(m),
		Data: m.GetData(),
	}
	if m.GetFileName() != "" {
		p.FileName = filepath.Base(m.GetFileName())
	}
	return json.Marshal(p)
}