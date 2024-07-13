package model

import (
	"fmt"
	"reflect"
	"strings"
)

var models = map[string]any{}

type GetFile interface {
	GetFile() error
}

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

func RegisterModel(model Data) {
	models[GetName(model)] = model
}

func GetNewModel(name string) (Data, error) {
	if model, ok := models[name]; ok {
		v := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		return v.(Data), nil
	}
	return nil, fmt.Errorf("model not found: %s", name)
}

func GetName(m any) string {
	p := strings.Split(reflect.TypeOf(m).Elem().PkgPath(), "/")
	return p[len(p)-1]
}
