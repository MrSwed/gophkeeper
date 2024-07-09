package input

import (
	"reflect"
	"strings"
	"testing"
)

func TestModel_Bytes(t *testing.T) {
	tests := []struct {
		name    string
		m       Model
		want    []byte
		wantErr bool
	}{
		{
			name: "test auth",
			m: &Auth{
				Data: AuthData{
					Login:    "test",
					Password: "password"},
			},
			want: []byte("test:password"),
		},
		{
			name: "test bin",
			m: &Bin{
				Data: BinData{
					Bin: []byte("test:test"),
				},
			},
			want:    []byte("test:test"),
			wantErr: false,
		},
		{
			name: "test card 1",
			m: &Card{
				Data: CardData{
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
			m: &Card{
				Data: CardData{
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
			m: &Text{
				Common: Common{},
				Data: TextData{
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
		m           Model
		wantErrKeys []string
	}{
		{
			name: "test auth 1",
			m: &Auth{
				Common: Common{
					Key: "somesite.com",
				},
				Data: AuthData{
					Login:    "test",
					Password: "password",
				},
			},
		},
		{
			name: "test auth 2, need key",
			m: &Auth{
				Data: AuthData{
					Login:    "test",
					Password: "password",
				},
			},
			wantErrKeys: []string{"Key"},
		},

		{
			name: "test bin",
			m: &Bin{
				Common: Common{
					Key: "some bin data",
				},
				Data: BinData{
					Bin: []byte("test:test"),
				},
			},
		},

		{
			name: "no card number",
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
					Exp:  "05/25",
					Name: "Some Name",
					CVV:  "999",
				},
			},
			wantErrKeys: []string{"Number"},
		},
		{
			name: "not valid card",
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
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
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
					Exp:    "05/25",
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},
		{
			name: "bad cvv 1",
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
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
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
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
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
					Exp:    "05/25",
					Number: "4012888888881881",
				},
			},
		},
		{
			name: "bad exp",
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
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
			m:           &Card{},
			wantErrKeys: []string{"Number", "Key"},
		},
		{
			name: "can no exp",
			m: &Card{
				Common: Common{Key: "some bank card 1"},
				Data: CardData{
					Number: "4012888888881881",
					Name:   "Some Name",
					CVV:    "999",
				},
			},
		},

		{
			name: "text 1",
			m: &Text{
				Common: Common{
					Key: "some test record 1",
				},
				Data: TextData{
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
