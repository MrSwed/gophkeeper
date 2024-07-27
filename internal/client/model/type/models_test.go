package _type

import (
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"
	"reflect"
	"strings"
	"testing"
)

func TestModel_Bytes(t *testing.T) {
	tests := []struct {
		name    string
		m       model.Model
		want    []byte
		wantErr bool
	}{
		{
			name: "test auth",
			m: &auth.Model{
				Data: &auth.Data{
					Login:    "test",
					Password: "password"},
			},
			want: []byte(`{"type":"auth","data":{"login":"test","password":"password"}}`),
		},
		{
			name: "test bin",
			m: &bin.Model{
				Data: &bin.Data{
					Bin: []byte("test:test"),
				},
			},
			want:    []byte(`{"type":"bin","data":{"bin":"dGVzdDp0ZXN0"}}`),
			wantErr: false,
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
			want:    []byte(`{"type":"card","data":{"exp":"05/25","number":"0000000000000000","name":"Some Name","cvv":"999"}}`),
			wantErr: false,
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
			want:    []byte(`{"type":"card","data":{"exp":"05/25","number":"0000000000000000","cvv":"999"}}`),
			wantErr: false,
		},
		{
			name: "test text ",
			m: &text.Model{
				Common: model.Common{},
				Data: &text.Data{
					Text: "some text here",
				},
			},
			want:    []byte(`{"type":"text","data":{"text":"some text here"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_Validate(t *testing.T) {
	tests := []struct {
		name        string
		m           model.Model
		fields      []string
		wantErrKeys []string
	}{
		{
			name: "test auth 1",
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
			name: "test auth 2, need key",
			m: &auth.Model{
				Data: &auth.Data{
					Login:    "test",
					Password: "password",
				},
			},
			wantErrKeys: []string{"Key"},
		},

		{
			name: "test bin",
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
			name: "no card number",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 1"},
				Data: &card.Data{
					Exp:  "05/25",
					Name: "Some Name",
					CVV:  "999",
				},
			},
			wantErrKeys: []string{"Number"},
		},
		{
			name: "not valid card",
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data: &card.Data{
					Exp:    "05/25",
					Number: "0000000000000001",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
			wantErrKeys: []string{"Number"},
		},
		{
			name: "valid data",
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
			name:   "card validate num only",
			fields: []string{"Number"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data: &card.Data{
					Number: "4012888888881881",
				},
			},
		},
		{
			name:   "bad card num",
			fields: []string{"Data.Number"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 2"},
				Data:   &card.Data{},
			},
			wantErrKeys: []string{"Data.Number"},
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
			wantErrKeys: []string{"CVV"},
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
			wantErrKeys: []string{"CVV"},
		},
		{
			name:   "bad cvv 3, partial validate",
			fields: []string{"Data.CVV"},
			m: &card.Model{
				Common: model.Common{Key: "some bank card 6"},
				Data: &card.Data{
					CVV: "99",
				},
			},
			wantErrKeys: []string{"Data.CVV"},
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
			wantErrKeys: []string{"Exp"},
		},
		{
			name:        "no data",
			m:           &card.Model{},
			wantErrKeys: []string{"Data", "Key"},
		},
		{
			name:        "bad data",
			m:           &card.Model{Data: &card.Data{}},
			wantErrKeys: []string{"Number", "Key"},
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
			err := tt.m.Validate(tt.fields...)
			if (err != nil) != (tt.wantErrKeys != nil) || (tt.wantErrKeys != nil) != containStrInErr(err, tt.wantErrKeys...) {
				t.Errorf("Validate() error = %v, wantErrKeys %v", err, tt.wantErrKeys)
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

func TestModel_GetNewDataModel(t *testing.T) {
	tests := []struct {
		name    string
		m       model.Data
		wantErr bool
	}{
		{
			name: "test auth",
			m:    &auth.Data{},
		},
		{
			name: "test text",
			m:    &text.Data{},
		},
		{
			name: "test card",
			m:    &card.Data{},
		},
		{
			name: "test bin",
			m:    &bin.Data{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.GetNewDataModel(model.GetName(tt.m))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNewDataModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.m) {
				t.Errorf("Type() got = %v, want %v", got, tt.m)
			}
		})
	}
}
