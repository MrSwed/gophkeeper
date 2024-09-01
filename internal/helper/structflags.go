package helper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

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
func GenerateFlags(dst interface{}, fs *pflag.FlagSet) (err error) {
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
		if tagNames[0] == "" {
			continue
		}
		usage := sf.Tag.Get("usage")
		defVal := sf.Tag.Get("default")

		switch fv.Kind() {
		case reflect.Struct:
			return GenerateFlags(fv, fs)
		case reflect.Bool, reflect.Int64, reflect.Float64, reflect.Int, reflect.Uint64, reflect.String:
			switch p := fv.Addr().Interface().(type) {
			case *bool:
				defVal, _ := strconv.ParseBool(defVal)
				fs.BoolVarP(p, tagNames[0], tagNames[1], defVal, usage)
			case *string:
				fs.StringVarP(p, tagNames[0], tagNames[1], defVal, usage)
			case *int:
				defVal, _ := strconv.Atoi(defVal)
				fs.IntVarP(p, tagNames[0], tagNames[1], defVal, usage)
			case *time.Duration:
				defVal, _ := time.ParseDuration(defVal)
				fs.DurationVarP(p, tagNames[0], tagNames[1], defVal, usage)
			case *float64:
				defVal, _ := strconv.ParseFloat(defVal, 64)
				fs.Float64VarP(p, tagNames[0], tagNames[1], defVal, usage)
			case *uint64:
				defVal, _ := strconv.ParseUint(defVal, 10, 64)
				fs.Uint64VarP(p, tagNames[0], tagNames[1], defVal, usage)
			}
		default:
			err = errors.New("unknown type")
		}
	}
	return
}
