package card

import (
	"gophKeeper/internal/client/model"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	expDateRegexp = `^(0[1-9]|1[0-2])[|/]?([0-9]{4}|[0-9]{2})$`
	cvvRegexp     = `^\d{3}$`
)

type Data struct {
	Exp    string `json:"exp" validate:"omitempty,credit_card_exp_date"`
	Number string `json:"number" validate:"required,credit_card"`
	Name   string `json:"name,omitempty" validate:"omitempty"`
	CVV    string `json:"cvv,omitempty" validate:"omitempty,credit_card_cvv"`
}

type Model struct {
	model.Common
	Data *Data `json:"data" validate:"required"`
}

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

func (m *Model) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m)
}

func (m *Model) GetData() any {
	return m.Data.GetData()
}

func (m *Data) GetData() any {
	return m
}

func init() {
	model.RegisterModel(&Data{})
	validators := map[string]string{
		"credit_card_exp_date": expDateRegexp,
		"credit_card_cvv":      cvvRegexp,
	}

	for k, v := range validators {
		err := model.Validator.RegisterValidation(k, func(fl validator.FieldLevel) bool {
			result, err := regexp.Match(v, []byte(fl.Field().String()))
			return result && err == nil
		})
		if err != nil {
			panic(err)
		}
	}
}
