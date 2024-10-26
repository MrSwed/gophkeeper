package crypt

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
			gotCipherText2, err := Encode(tt.args.plainText, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, false, reflect.DeepEqual(gotCipherText, gotCipherText2),
				"twice ciphered can not be equal")

			gotPlainText, err := Decode(gotCipherText, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotPlainText2, err := Decode(gotCipherText2, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, true, reflect.DeepEqual(gotPlainText, gotPlainText2),
				fmt.Sprintf("Decode() gotPlainText = %v, gotPlainText2 %v", gotPlainText, gotPlainText2))

			require.Equal(t, true, reflect.DeepEqual(gotPlainText, tt.args.plainText),
				fmt.Sprintf("Decode() gotPlainText = %v, want %v", gotPlainText, tt.args.plainText))

			require.Equal(t, false, reflect.DeepEqual(gotPlainText, gotCipherText),
				fmt.Sprintf("Decode() gotPlainText = %v, gotCipherText %v", gotPlainText, gotCipherText))
		})
	}
}
