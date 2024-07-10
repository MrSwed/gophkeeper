package models

import (
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/models/auth"
	"gophKeeper/internal/client/model/models/bin"
	"gophKeeper/internal/client/model/models/card"
	"gophKeeper/internal/client/model/models/text"
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
			m: &auth.Auth{
				Data: auth.AuthData{
					Login:    "test",
					Password: "password"},
			},
			want: []byte("test:password"),
		},
		{
			name: "test bin",
			m: &bin.Bin{
				Data: bin.BinData{
					Bin: []byte("test:test"),
				},
			},
			want:    []byte("test:test"),
			wantErr: false,
		},
		{
			name: "test card 1",
			m: &card.Card{
				Data: card.CardData{
					Exp:    "05/25",
					Number: "0000000000000000",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
			want:    []byte("0000000000000000|05/25|999|Some Name"),
			wantErr: false,
		},
		{
			name: "test card 2",
			m: &card.Card{
				Data: card.CardData{
					Exp:    "05/25",
					Number: "0000000000000000",
					CVV:    "999",
				},
			},
			want:    []byte("0000000000000000|05/25|999|"),
			wantErr: false,
		},
		{
			name: "test text ",
			m: &text.Text{
				Common: model.Common{},
				Data: text.TextData{
					Text: "some text here",
				},
			},
			want:    []byte("some text here"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Bytes()
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
		wantErrKeys []string
	}{
		{
			name: "test auth 1",
			m: &auth.Auth{
				Common: model.Common{
					Key: "somesite.com",
				},
				Data: auth.AuthData{
					Login:    "test",
					Password: "password",
				},
			},
		},
		{
			name: "test auth 2, need key",
			m: &auth.Auth{
				Data: auth.AuthData{
					Login:    "test",
					Password: "password",
				},
			},
			wantErrKeys: []string{"Key"},
		},

		{
			name: "test bin",
			m: &bin.Bin{
				Common: model.Common{
					Key: "some bin data",
				},
				Data: bin.BinData{
					Bin: []byte("test:test"),
				},
			},
		},

		{
			name: "no card number",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
					Exp:  "05/25",
					Name: "Some Name",
					CVV:  "999",
				},
			},
			wantErrKeys: []string{"Number"},
		},
		{
			name: "not valid card",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
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
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},
		{
			name: "bad cvv 1",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
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
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "99",
				},
			},
			wantErrKeys: []string{"CVV"},
		},
		{
			name: "can be no cvv",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
					Exp:    "05/25",
					Number: "4012888888881881",
				},
			},
		},
		{
			name: "bad exp",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
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
			m:           &card.Card{},
			wantErrKeys: []string{"Number", "Key"},
		},
		{
			name: "can no exp",
			m: &card.Card{
				Common: model.Common{Key: "some bank card 1"},
				Data: card.CardData{
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},

		{
			name: "text 1",
			m: &text.Text{
				Common: model.Common{
					Key: "some test record 1",
				},
				Data: text.TextData{
					Text: "some text here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
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
