package model

import (
	"fmt"
	"reflect"
	"strings"
)

var models = map[string]any{}

type Data interface {
	GetData() any
}

type Base interface {
	GetKey() string
	GetDescription() *string
	GetFileName() string
}

type Model interface {
	Validate
	Base
	Bytes() (b []byte, err error)
	Data
}

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

func RegisterModel(model Data) {
	models[GetName(model)] = model
}

func GetNewModel(name string) (Data, error) {
	if model, ok := models[name]; ok {
		v := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		return v.(Data), nil
	}
	return nil, fmt.Errorf("model %s not found", name)
}

func GetName(m any) string {
	p := strings.Split(reflect.TypeOf(m).Elem().PkgPath(), "/")
	return p[len(p)-1]
}
