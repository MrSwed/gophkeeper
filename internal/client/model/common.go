package model

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"
)

type Common struct {
	Key         string `json:"key" validate:"required" flag:"key,k" usage:"set your entry key-identifier"`
	Description string `json:"description" flag:"file,f" usage:"read from file"`
	FileName    string `json:"fileName" flag:"description,d" usage:"description, will be displayed in the list of entries list"`
}

func (c *Common) GetKey() string {
	if c.Key == "" {
		c.Key = fmt.Sprintf("Key-%s", time.Now().Format("2006-01-02-15-04-05"))
	}
	return c.Key
}

func (c *Common) GetDescription() string {
	return c.Description
}
func (c *Common) GetFileName() string {
	return c.FileName
}

func (c *Common) GetBase() *Common {
	return c
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
		Data: m.GetPacked(),
	}
	if m.GetFileName() != "" {
		p.FileName = filepath.Base(m.GetFileName())
	}
	return json.Marshal(p)
}
