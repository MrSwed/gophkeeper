package data

import (
	"reflect"
	"testing"
)

func TestAuth_Bytes(t *testing.T) {
	type fields struct {
		Common   Common
		Login    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Login:    "test",
				Password: "password",
			},
			want: []byte("test:password"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Auth{
				Common:   tt.fields.Common,
				Login:    tt.fields.Login,
				Password: tt.fields.Password,
			}
			got, err := m.Bytes()
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

func TestAuth_IsValid(t *testing.T) {
	type fields struct {
		Common   Common
		Login    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Common: Common{
					Key: "somesite.com",
				},
				Login:    "test",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "test 2, need key",
			fields: fields{
				Login:    "test",
				Password: "password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Auth{
				Common:   tt.fields.Common,
				Login:    tt.fields.Login,
				Password: tt.fields.Password,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBin_Bytes(t *testing.T) {
	type fields struct {
		Common Common
		Bin    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Bin: []byte("test:test"),
			},
			want:    []byte("test:test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Bin{
				Common: tt.fields.Common,
				Bin:    tt.fields.Bin,
			}
			got, err := m.Bytes()
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

func TestBin_IsValid(t *testing.T) {
	type fields struct {
		Common Common
		Bin    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Common: Common{
					Key: "some bin data",
				},
				Bin: []byte("test:test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Bin{
				Common: tt.fields.Common,
				Bin:    tt.fields.Bin,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCard_Bytes(t *testing.T) {
	type fields struct {
		Common Common
		Exp    string
		Number string
		Name   string
		CVV    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Exp:    "05/25",
				Number: "0000000000000000",
				Name:   "Some Name",
				CVV:    "999",
			},
			want:    []byte("0000000000000000|05/25|999|Some Name"),
			wantErr: false,
		},
		{
			name: "test 2",
			fields: fields{
				Exp:    "05/25",
				Number: "0000000000000000",
				CVV:    "999",
			},
			want:    []byte("0000000000000000|05/25|999|"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Card{
				Common: tt.fields.Common,
				Exp:    tt.fields.Exp,
				Number: tt.fields.Number,
				Name:   tt.fields.Name,
				CVV:    tt.fields.CVV,
			}
			got, err := m.Bytes()
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

func TestCard_IsValid(t *testing.T) {
	type fields struct {
		Common Common
		Exp    string
		Number string
		Name   string
		CVV    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "no card",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Name:   "Some Name",
				CVV:    "999",
			},
			wantErr: true,
		},
		{
			name: "not valid card",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Number: "0000000000000001",
				Name:   "Some Name",
				CVV:    "999",
			},
			wantErr: true,
		},
		{
			name: "valid data",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Number: "4012888888881881",
				Name:   "Some Name",
				CVV:    "999",
			},
			wantErr: false,
		},
		{
			name: "bad cvv 1",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Number: "4012888888881881",
				Name:   "Some Name",
				CVV:    "9992",
			},
			wantErr: true,
		},
		{
			name: "bad cvv 2",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Number: "4012888888881881",
				Name:   "Some Name",
				CVV:    "99",
			},
			wantErr: true,
		},
		{
			name: "can be no cvv",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "05/25",
				Number: "4012888888881881",
			},
			wantErr: false,
		},
		{
			name: "bad exp",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Exp:    "48/25",
				Number: "4012888888881881",
				Name:   "Some Name",
				CVV:    "999",
			},
			wantErr: true,
		},
		{
			name: "can no exp",
			fields: fields{
				Common: Common{Key: "some bank card 1"},
				Number: "4012888888881881",
				Name:   "Some Name",
				CVV:    "999",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Card{
				Common: tt.fields.Common,
				Exp:    tt.fields.Exp,
				Number: tt.fields.Number,
				Name:   tt.fields.Name,
				CVV:    tt.fields.CVV,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestText_Bytes(t *testing.T) {
	type fields struct {
		Common Common
		Text   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test 1 ",
			fields: fields{
				Common: Common{},
				Text:   "some text here",
			},
			want:    []byte("some text here"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Text{
				Common: tt.fields.Common,
				Text:   tt.fields.Text,
			}
			got, err := m.Bytes()
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

func TestText_IsValid(t *testing.T) {
	type fields struct {
		Common Common
		Text   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				Common: Common{
					Key: "some test record 1",
				},
				Text: "some text here",
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Text{
				Common: tt.fields.Common,
				Text:   tt.fields.Text,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
