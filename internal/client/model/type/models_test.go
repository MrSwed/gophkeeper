package _type

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ model.Model = (*unkModel)(nil)
	// _ model.Data  = (*Data)(nil)
)

type unkModel struct {
	model.Data
	model.Common
}

func (m *unkModel) GetKey() string         { return "" }
func (m *unkModel) GetDescription() string { return "" }
func (m *unkModel) GetBase() *model.Common { return &m.Common }
func (m *unkModel) Reset() {
	m.Common.Reset()
	m.Data.Reset()
}
func (m *unkModel) Validate(_ ...string) error { return nil }
func (m *unkModel) GetPacked() any             { return m.Data.GetPacked() }
func (m *unkModel) GetDst() any                { return m.Data.GetDst() }

func TestModel(t *testing.T) {
	tests := []struct {
		name                string
		wantBytes           []byte
		validate            []string
		validateWantErrKeys []string
		m                   model.Model
		detectModelErr      bool
		packedErr           bool
	}{
		{
			name: "auth",
			m: &auth.Model{
				Data: &auth.Data{
					Login:    "test",
					Password: "password"},
			},
			wantBytes: []byte(`{"type":"auth","data":{"login":"test","password":"password"}}`),
		},
		{
			name: "bin",
			m: &bin.Model{
				Data: &bin.Data{
					Bin: []byte("test:test"),
				},
			},
			wantBytes: []byte(`{"type":"bin","data":{"bin":"dGVzdDp0ZXN0"}}`),
		},
		{
			name: "card 1",
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
			name: "card 2",
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
			name: "text",
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
			name:     "auth 1",
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
			name:     "auth 2, need key",
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
			name:     "bin",
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
			name:                "new auth",
			m:                   auth.New(),
			validate:            []string{},
			validateWantErrKeys: []string{"Key"},
		},
		{
			name:                "no card",
			m:                   card.New(),
			validate:            []string{},
			validateWantErrKeys: []string{"Key"},
		},
		{
			name:                "no text",
			m:                   text.New(),
			validate:            []string{},
			validateWantErrKeys: []string{"Key"},
		},
		{
			name:                "no bin",
			m:                   bin.New(),
			validate:            []string{},
			validateWantErrKeys: []string{"Key"},
		},
		{
			name:     "card can no exp",
			validate: []string{},
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
			name:     "text 1",
			validate: []string{},
			m: &text.Model{
				Common: model.Common{
					Key: "some record 1",
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
					if (err != nil) != tt.packedErr {
						t.Errorf("PackedBytes() error = %v, wantErr %v", err, tt.detectModelErr)
						return
					}
					if !reflect.DeepEqual(got, tt.wantBytes) {
						t.Errorf("PackedBytes() got = %s, want %s, data %v", string(got), string(tt.wantBytes), tt.m.GetDst())
					}
				})
			}

			t.Run("GetBase", func(t *testing.T) {
				got := tt.m.GetBase()
				assert.IsType(t, &model.Common{}, got, fmt.Errorf("GetBase() got = %v, for model %v", got, tt.m))
			})

			if !tt.detectModelErr {
				t.Run("GetKey", func(t *testing.T) {
					oldKey := tt.m.GetBase().Key
					updatedKey := tt.m.GetKey()
					require.Equal(t, true, updatedKey != "", "GetKey should not be empty")
					if oldKey != "" {
						require.Equal(t, oldKey, updatedKey, "GetKey key should not be changed if it not empty already")
					}
					require.Equal(t, true, tt.m.GetBase().GetKey() != "", "Base GetKey also should not be empty")
				})

				t.Run("GetDescription", func(t *testing.T) {
					require.Equal(t, tt.m.GetBase().Description, tt.m.GetDescription(), "GetDescription should just get description")
				})

				t.Run("Reset", func(t *testing.T) {

					nModel := &unkModel{}

					var err error
					nModel.Data, err = model.GetNewDataModel(model.GetName(tt.m))
					require.NoError(t, err)

					tt.m.Reset()

					require.Equal(t, true, reflect.DeepEqual(nModel.GetDst(), tt.m.GetDst()),
						fmt.Errorf("not empty GetDst() :  new = %v, tt.m %v", nModel.GetDst(), tt.m.GetDst()))

					require.Equal(t, true, reflect.DeepEqual(nModel.GetBase(), tt.m.GetBase()),
						fmt.Errorf("not empty GetBase() :  new = %v, tt.m %v", nModel.GetBase(), tt.m.GetBase()))

				})
			} else if tt.m.GetBase() != nil {
				t.Run("GetKey", func(t *testing.T) {
					require.Equal(t, true, tt.m.GetBase().GetKey() != "", "Base GetKey also should not be empty")
				})
			}
		})
	}
}

func TestGetNewDataModel(t *testing.T) {
	tests := []struct {
		name    string
		want    model.Data
		wantErr error
	}{
		{
			name: "auth",
			want: &auth.Data{},
		},
		{
			name: "text",
			want: &text.Data{},
		},
		{
			name: "bin",
			want: &bin.Data{},
		},
		{
			name: "card",
			want: &card.Data{},
		},
		{
			name:    "unknown",
			wantErr: fmt.Errorf("model not found: %s", "unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.GetNewDataModel(tt.name)
			require.Equal(t, true,
				(err == nil && tt.wantErr == nil) ||
					errors.Is(err, tt.wantErr) ||
					err.Error() == tt.wantErr.Error(),
				tt.name,
				fmt.Errorf("GetNewDataModel() error = %v, wantErr %v", err, tt.wantErr))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNewDataModel() got = %v, want %v", got, tt.want)
			}
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
