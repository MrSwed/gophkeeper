package model

import (
	"fmt"
	"reflect"
	"strings"
)

var models = map[string]any{}

func RegisterModel(model Model) {
	models[model.Type()] = model
}

func GetNewModel(name string) (any, error) {
	if model, ok := models[name]; ok {
		v := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		return v, nil
	}
	return nil, fmt.Errorf("model %s not found", name)
}

func GetName(m any) string {
	p := strings.Split(reflect.TypeOf(m).Elem().PkgPath(), "/")
	return p[len(p)-1]
}
