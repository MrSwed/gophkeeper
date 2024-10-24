package helper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	StringVar string  `json:"string" flag:"string,s" default:"" usage:"usage for stringVar"`
	FloatVar  float64 `json:"float" flag:"float,f" default:"" usage:"usage for floatVar"`
	UintVar   uint64  `json:"uint" flag:"uint,u" default:"" usage:"usage for uintVar"`
	IntVar    int     `json:"int" flag:"int,i" default:"" usage:"usage for intVar"`
}

func (ts *testStruct) check(t *testing.T, _ map[string]string) {
	require.Equal(t, true, ts.StringVar != "0")
	require.Equal(t, true, ts.IntVar != 0)
	require.Equal(t, true, ts.FloatVar != 0)
	require.Equal(t, true, ts.UintVar != 0)
}

type testErrStruct struct {
	BadVar    map[string]any `json:"bad" flag:"bad,b" default:"" usage:"usage for badVar"`
	StringVar string         `json:"string" flag:"string,s" default:"" usage:"usage for stringVar"`
}

type testStructSub struct {
	DurationVar time.Duration `json:"duration" flag:"duration,d" default:"" usage:"usage for durationVar"`
	BoolVar     bool          `json:"bool" flag:"bool,b" default:"" usage:"usage for boolVar"`
}

func (ts *testStructSub) check(t *testing.T, want map[string]string) {
	require.Equal(t, true, ts.DurationVar != 0)
	var (
		b   bool
		err error
	)
	if s, ok := want["bool"]; ok {
		b, err = strconv.ParseBool(s)
		require.NoError(t, err)
	}
	require.Equal(t, true, ts.BoolVar == b)
}

type testStructCollect struct {
	testStruct
	testStructSub
}

func TestGenerateFlags(t *testing.T) {

	type args struct {
		dst any
		fs  *pflag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    map[string]string
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
			name: "test struct sub",
			args: args{
				dst: &testStructCollect{},
				fs:  (&cobra.Command{}).Flags(),
			},
			want: map[string]string{
				"string":   "string",
				"int":      "-1",
				"float":    "0.1",
				"uint":     "2",
				"bool":     "true",
				"duration": "1m0s",
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

			if tt.wantErr == nil {
				for k, v := range tt.want {
					err = tt.args.fs.Set(k, v)
					require.NoError(t, err)
				}
				if dst, ok := tt.args.dst.(*testStruct); ok {
					dst.check(t, nil)
				}
				if dst, ok := tt.args.dst.(*testStructSub); ok {
					dst.check(t, tt.want)
				}
				if dst, ok := tt.args.dst.(*testStructCollect); ok {
					dst.testStruct.check(t, nil)
					dst.testStructSub.check(t, tt.want)
				}
			}
		})
	}
}
