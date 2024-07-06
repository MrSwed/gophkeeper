package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	expDateRegexp = `^(0[1-9]|1[0-2])[|/]?([0-9]{4}|[0-9]{2})$`
	cvvRegexp     = `^\d{3}$`
)

type Card struct {
	Common
	Exp    string `json:"exp" validate:"omitempty,credit_card_exp_date"`
	Number string `json:"number" validate:"required,credit_card"`
	Name   string `json:"name" validate:"omitempty"`
	CVV    string `json:"cvv" validate:"omitempty,credit_card_cvv"`
}

func (m *Card) Validate() error {
	return validate.Struct(m)
}

func (m *Card) Bytes() ([]byte, error) {
	return []byte(fmt.Sprintf("%s|%s|%s|%s", m.Number, m.Exp, m.CVV, m.Name)), nil
}

func init() {
	validators := map[string]string{
		"credit_card_exp_date": expDateRegexp,
		"credit_card_cvv":      cvvRegexp,
	}

	for k, v := range validators {
		err := validate.RegisterValidation(k, func(fl validator.FieldLevel) bool {
			result, err := regexp.Match(v, []byte(fl.Field().String()))
			return result && err == nil
		})
		if err != nil {
			panic(err)
		}
	}
}