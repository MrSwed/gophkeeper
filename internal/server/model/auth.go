package model

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Meta     any    `json:"-"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (m *AuthRequest) Validate(fields ...string) error {
	return ValidateStruct(m, fields...)
}

func init() {
	validators := map[string][]string{
		"password": {".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"},
	}

	for k, vv := range validators {
		err := Validator.RegisterValidation(k, func(fl validator.FieldLevel) bool {
			for _, v := range vv {
				result, err := regexp.Match(v, []byte(fl.Field().String()))
				if !result || err != nil {
					return false
				}
			}
			return true
		})
		if err != nil {
			panic(err)
		}
	}
}
