package helper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

// GenerateFlags
//
//		set pflag.FlagSet for structure dst
//	 struct tags used for pflag params
//	 - flag - long,short flags (separate comma)
//	 - default - will be converted to it type, see supported types in code
//	 - usage - usage text
//
// thanks https://stackoverflow.com/questions/72891199/procedurally-bind-struct-fields-to-command-line-flag-values-using-reflect
func GenerateFlags(dst interface{}, fs *pflag.FlagSet) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("not pointer-to-a-struct") // exit if not pointer-to-a-struct
	}

	rv = rv.Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)
		tagNames := [2]string{}
		copy(tagNames[:], strings.SplitN(sf.Tag.Get(("flag")), ",", 2))
		usage := sf.Tag.Get("usage")
		defVal := sf.Tag.Get("default")

		switch fv.Type() {
		case reflect.TypeOf(string("")):
			p := fv.Addr().Interface().(*string)
			fs.StringVarP(p, tagNames[0], tagNames[1], defVal, usage)
		case reflect.TypeOf(int(0)):
			p := fv.Addr().Interface().(*int)
			defVal, _ := strconv.Atoi(defVal)
			fs.IntVarP(p, tagNames[0], tagNames[1], defVal, usage)
		case reflect.TypeOf(float64(0)):
			p := fv.Addr().Interface().(*float64)
			defVal, _ := strconv.ParseFloat(defVal, 64)
			fs.Float64VarP(p, tagNames[0], tagNames[1], defVal, usage)
		case reflect.TypeOf(uint64(0)):
			p := fv.Addr().Interface().(*uint64)
			defVal, _ := strconv.ParseUint(defVal, 10, 64)
			fs.Uint64VarP(p, tagNames[0], tagNames[1], defVal, usage)
		default:
			return GenerateFlags(fv.Interface(), fs)
		}
	}
	return nil
}
