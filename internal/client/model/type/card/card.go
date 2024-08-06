package card

import (
	"encoding/json"
	"fmt"
	"gophKeeper/internal/client/model"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	expDateRegexp = `^(0[1-9]|1[0-2])[|/]?([0-9]{4}|[0-9]{2})$`
	cvvRegexp     = `^\d{3}$`
)

var (
	_ model.Model = (*Model)(nil)
	_ model.Data  = (*Data)(nil)
)

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

type Model struct {
	model.Common
	Data *Data `json:"data" validate:"required"`
}

func (m *Model) GetKey() string {
	if m.Key == "" {
		m.Key = fmt.Sprintf("%s-%s", model.GetName(m), time.Now().Format("2006-01-02-15-04-05"))
	}
	return m.Key
}

func (m *Model) Validate(fields ...string) error {
	if len(fields) == 0 {
		return model.Validator.Struct(m)
	} else {
		return model.Validator.StructPartial(m, fields...)
	}
}

func (m *Model) Bytes() (b []byte, err error) {
	return model.NewPackedBytes(m)
}

func (m *Model) GetData() any {
	return m.Data.GetData()
}

type Data struct {
	Exp    string `json:"exp" validate:"omitempty,credit_card_exp_date"`
	Number string `json:"number" validate:"required,credit_card"`
	CVV    string `json:"cvv,omitempty" validate:"omitempty,credit_card_cvv"`
	Name   string `json:"name,omitempty" validate:"omitempty"`
}

type packedData struct {
	Number cardNumber `json:"number"`
	Exp    cardExo    `json:"exp"`
	CVV    string     `json:"cvv,omitempty"`
	Name   string     `json:"name,omitempty"`
}

func (m *Data) GetData() any {
	p := new(packedData)
	p.CVV = m.CVV
	p.Name = m.Name
	p.Exp.Set(m.Exp)
	p.Number.Set(m.Number)
	return p
}

func (m *Data) Sanitize() {
	if packed, _ := m.GetData().(*packedData); packed != nil {
		m.Exp = packed.Exp.String()
		m.Number = packed.Number.String()
	}
}

type cardNumber [16]byte

func (c *cardNumber) Set(s string) {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	*c = cardNumber{}
	for i := 0; i < len(s) && i < len(c); i++ {
		c[i] = s[i]
	}
}

func (c *cardNumber) String() string {
	b := make([]byte, 0, len(c))
	for i := 0; i < len(c); i++ {
		if c[i] != 0 {
			b = append(b, c[i])
			if i > 0 && (i+1)%4 == 0 {
				b = append(b, ' ')
			}
		}
	}

	return strings.TrimSpace(string(b))
}

func (c *cardNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(string((*c)[:]))
}

type cardExo [4]byte

func (c *cardExo) Set(s string) {
	s = strings.ReplaceAll(s, "/", "")
	*c = cardExo{}
	copy(c[:], s)
}

func (c *cardExo) String() string {
	b := make([]byte, 0, len(c))
	for i := 0; i < len(c); i++ {
		if c[i] != 0 {
			b = append(b, c[i])
		}
	}
	if len(b) > 2 {
		b = append([]byte{b[0], b[1], '/'}, b[2:]...)
	}
	return strings.TrimSpace(string(b))
}

func (c *cardExo) MarshalJSON() ([]byte, error) {
	return json.Marshal(string((*c)[:]))
}
