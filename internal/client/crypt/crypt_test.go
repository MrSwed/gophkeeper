package crypt

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	type args struct {
		plainText []byte
		key       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				plainText: []byte("some sext1"),
				key:       "someKeyPhraseSecret",
			},
		},
		{
			name: "test 2",
			args: args{
				plainText: []byte("00001111222233330519333"),
				key:       "somesecretkey2",
			},
		},
		{
			name: "test empty key",
			args: args{
				plainText: []byte("00001111222233330519333"),
				key:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCipherText, err := Encode(tt.args.plainText, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotPlainText, err := Decode(gotCipherText, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotPlainText, tt.args.plainText) {
				t.Errorf("Decode() gotPlainText = %v, want %v", gotPlainText, tt.args.plainText)
			}

			if reflect.DeepEqual(gotPlainText, gotCipherText) {
				t.Errorf("Decode() gotPlainText = %v, gotCipherText %v", gotPlainText, gotCipherText)
			}
		})
	}
}
