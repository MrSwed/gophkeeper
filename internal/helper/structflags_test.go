package helper

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestGenerateFlags(t *testing.T) {

	type testStruct struct {
		StringVar string  `json:"string" flag:"string,s" default:"" usage:"usage for stringVar"`
		IntVar    int     `json:"int" flag:"int,i" default:"" usage:"usage for intVar"`
		FloatVar  float64 `json:"float" flag:"float,f" default:"" usage:"usage for floatVar"`
		UintVar   uint64  `json:"uint" flag:"uint,u" default:"" usage:"usage for uintVar"`
	}

	type testErrStruct struct {
		StringVar string         `json:"string" flag:"string,s" default:"" usage:"usage for stringVar"`
		BadVar    map[string]any `json:"bad" flag:"bad,b" default:"" usage:"usage for badVar"`
	}

	type args struct {
		dst any
		fs  *pflag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr error
	}{
		{
			name: "test ok",
			args: args{
				dst: &testStruct{},
				fs:  (&cobra.Command{}).Flags(),
			},
			want: map[string]string{
				"string": "string",
				"int":    "-1",
				"float":  "0.1",
				"uint":   "2",
			},
		},
		{
			name: "test err",
			args: args{
				dst: &testErrStruct{},
				fs:  (&cobra.Command{}).Flags(),
			},
			wantErr: errors.New("unknown type"),
		},
		{
			name: "test err 2",
			args: args{
				dst: "not pointer-to-a-struct",
				fs:  (&cobra.Command{}).Flags(),
			},
			wantErr: errors.New("not pointer-to-a-struct"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateFlags(tt.args.dst, tt.args.fs)
			require.Equal(t, true,
				reflect.DeepEqual(err, tt.wantErr),
				fmt.Sprintf("GenerateFlags() error = %v, wantErr %v, data %v", err, tt.wantErr, tt.args.fs))

			if tt.wantErr != nil {
				for k, v := range tt.want {
					err = tt.args.fs.Set(k, v)
					require.NoError(t, err)
				}
				if dst, ok := tt.args.dst.(*testStruct); ok {
					require.Equal(t, true, dst.StringVar != "0")
					require.Equal(t, true, dst.IntVar != 0)
					require.Equal(t, true, dst.FloatVar != 0)
					require.Equal(t, true, dst.UintVar != 0)
				}
			}
		})
	}
}
