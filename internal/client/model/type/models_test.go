package _type

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_ model.Model = (*unkModel)(nil)
	// _ model.Data  = (*Data)(nil)
)

type unkModel struct {
	model.Common
	data []byte
}

func (m *unkModel) GetKey() string             { return "" }
func (m *unkModel) GetDescription() string     { return "" }
func (m *unkModel) GetBase() *model.Common     { return &m.Common }
func (m *unkModel) Reset()                     {}
func (m *unkModel) Validate(_ ...string) error { return nil }
func (m *unkModel) GetPacked() any             { return &m.data }
func (m *unkModel) GetDst() any                { return &m.data }

func TestModel(t *testing.T) {
	tests := []struct {
		name                string
		m                   model.Model
		validate            []string
		validateWantErrKeys []string
		detectModelErr      bool
		wantBytes           []byte
	}{
		{
			name: "test auth",
			m: &auth.Model{
				Data: &auth.Data{
					Login:    "test",
					Password: "password"},
			},
			wantBytes: []byte(`{"type":"auth","data":{"login":"test","password":"password"}}`),
		},
		{
			name: "test bin",
			m: &bin.Model{
				Data: &bin.Data{
					Bin: []byte("test:test"),
				},
			},
			wantBytes: []byte(`{"type":"bin","data":{"bin":"dGVzdDp0ZXN0"}}`),
		},
		{
			name: "test card 1",
			m: &card.Model{
				Data: &card.Data{
					Exp:    "05/25",
					Number: "0000000000000000",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
			wantBytes: []byte(`{"type":"card","data":{"number":"0000000000000000","exp":"0525","cvv":"999","name":"Some Name"}}`),
		},
		{
			name: "test card 2",
			m: &card.Model{
				Data: &card.Data{
					Exp:    "05/25",
					Number: "0000000000000000",
					CVV:    "999",
				},
			},
			wantBytes: []byte(`{"type":"card","data":{"number":"0000000000000000","exp":"0525","cvv":"999"}}`),
		},
		{
			name: "test text ",
			m: &text.Model{
				Common: model.Common{},
				Data: &text.Data{
					Text: "some text here",
				},
			},
			wantBytes: []byte(`{"type":"text","data":{"text":"some text here"}}`),
		},
		{
			name: "unknown model",
			m: &unkModel{
				Common: model.Common{},
			},
			detectModelErr: true,
		},
		{
			name:     "test auth 1",
			validate: []string{},
			m: &auth.Model{
				Common: model.Common{
					Key: "somesite.com",
				},
				Data: &auth.Data{
					Login:    "test",
					Password: "password",
				},
			},
		},
		{
			name:     "test auth 2, need key",
			validate: []string{},
			m: &auth.Model{
				Data: &auth.Data{
					Login:    "test",
					Password: "password",
				},
			},
			validateWantErrKeys: []string{"Key"},
		},

		{
			name:     "test bin",
			validate: []string{},
			m: &bin.Model{
				Common: model.Common{
					Key: "some bin data",
				},
				Data: &bin.Data{
					Bin: []byte("test:test"),
				},
			},
		},

		{
			name:     "no card number",
			validate: []string{},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 1"},
				Data: &card.Data{
					Exp:  "05/25",
					Name: "Some Name",
					CVV:  "999",
				},
			},
			validateWantErrKeys: []string{"Number"},
		},
		{
			name:     "not valid card",
			validate: []string{},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "0000000000000001",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
			validateWantErrKeys: []string{"Number"},
		},
		{
			name:     "valid data",
			validate: []string{},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 3"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},
		{
			name:     "card validate num only",
			validate: []string{"Number"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data: &card.Data{
					Number: "4012888888881881",
				},
			},
		},
		{
			name:     "bad card num",
			validate: []string{"Data.Number"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data:   &card.Data{},
			},
			validateWantErrKeys: []string{"Data.Number"},
		},
		{
			name: "bad cvv 1",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 4"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "9992",
				},
			},
			validateWantErrKeys: []string{"CVV"},
		},
		{
			name: "bad cvv 2",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 5"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "99",
				},
			},
			validateWantErrKeys: []string{"CVV"},
		},
		{
			name:     "bad cvv 3, partial validate",
			validate: []string{"Data.CVV"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 6"},
				Data: &card.Data{
					CVV: "99",
				},
			},
			validateWantErrKeys: []string{"Data.CVV"},
		},
		{
			name: "can be no cvv",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 7"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "4012888888881881",
				},
			},
		},
		{
			name: "bad exp",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 8"},
				Data: &card.Data{
					Exp:    "48/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
			validateWantErrKeys: []string{"Exp"},
		},
		{
			name:                "no data",
			m:                   card.New(),
			validateWantErrKeys: []string{"Data", "Key"},
		},
		{
			name:                "bad data",
			m:                   &card.Model{Data: &card.Data{}},
			validateWantErrKeys: []string{"Number", "Key"},
		},
		{
			name: "can no exp",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 9"},
				Data: &card.Data{
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},

		{
			name: "text 1",
			m: &text.Model{
				Common: model.Common{
					Key: "some test record 1",
				},
				Data: &text.Data{
					Text: "some text here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Run("Data model", func(t *testing.T) {
				_, err := model.GetNewDataModel(model.GetName(tt.m))
				if (err != nil) != tt.detectModelErr {
					t.Errorf("GetNewDataModel() error = %v, wantErr %v", err, tt.detectModelErr)
					return
				}
			})

			if tt.validate != nil {
				t.Run("Validate", func(t *testing.T) {
					err := tt.m.Validate(tt.validate...)
					if (err != nil) != (tt.validateWantErrKeys != nil) || (tt.validateWantErrKeys != nil) != containStrInErr(err, tt.validateWantErrKeys...) {
						t.Errorf("Validate() error = %v, wantErrKeys %v", err, tt.validateWantErrKeys)
					}
				})
			}

			if tt.wantBytes != nil {
				t.Run("PackedBytes", func(t *testing.T) {
					got, err := model.NewPackedBytes(tt.m)
					if (err != nil) != tt.detectModelErr {
						t.Errorf("PackedBytes() error = %v, wantErr %v", err, tt.detectModelErr)
						return
					}
					if !reflect.DeepEqual(got, tt.wantBytes) {
						t.Errorf("PackedBytes() got = %v, want %v", got, tt.wantBytes)
					}
				})
			}

			t.Run("GetBase", func(t *testing.T) {
				got := tt.m.GetBase()
				assert.IsType(t, &model.Common{}, got, fmt.Errorf("GetBase() got = %v, for model %v", got, tt.m))
			})

		})
	}
}

func containStrInErr(err error, str ...string) bool {
	if err == nil {
		return false
	}
	c := 0
	for _, s := range str {
		if strings.Contains(err.Error(), s) {
			c++
		}
	}
	return c == len(str)
}
