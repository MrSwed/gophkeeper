package card

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	expDateRegexp = `^(0[1-9]|1[0-2])[|/]?([0-9]{4}|[0-9]{2})$`
	cvvRegexp     = `^\d{3}$`
)

type ModelData struct {
	Exp    string `json:"exp" validate:"omitempty,credit_card_exp_date"`
	Number string `json:"number" validate:"required,credit_card"`
	Name   string `json:"name" validate:"omitempty"`
	CVV    string `json:"cvv" validate:"omitempty,credit_card_cvv"`
}

type Model struct {
	model.Common
	Data ModelData `json:"data"`
}

var _ model.Model = (*Model)(nil)

func (m *Model) Validate() error {
	return model.Validator.Struct(m)
}

func (m *Model) Bytes() []byte {
	return []byte(fmt.Sprintf("%s|%s|%s|%s", m.Data.Number, m.Data.Exp, m.Data.CVV, m.Data.Name))
}

func (m *Model) Type() string {
	return model.GetName(m)
}

func init() {
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
