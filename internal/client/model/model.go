package model

import (
	"fmt"
	"reflect"
	"strings"
)

var models = map[string]any{}

type FromFile interface {
	DataFromFile() error
}

type Data interface {
	GetPacked() any
	GetDst() any
	Reset()
}

type Base interface {
	GetKey() string
	GetDescription() string
	GetFileName() string
	GetBase() *Common
}

type Model interface {
	Validate
	Base
	Data
}

type Settable interface {
	Set(s string)
}

type Sanitisable interface {
	Sanitize()
}

func RegisterModel(model Data) {
	models[GetName(model)] = model
}

func GetNewDataModel(name string) (Data, error) {
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
