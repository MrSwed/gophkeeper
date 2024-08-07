package cmd

import (
	"errors"
	"fmt"
	cfg "gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/service"
	"gophKeeper/internal/client/storage"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App interface {
	Execute() error
}

type app struct {
	db   *sqlx.DB
	srv  service.Service
	root *cobra.Command
}

func NewApp() (a *app) {
	a = &app{}
	dbFile := cfg.User.GetString("db_file")
	if dbFile == "" {
		err := errors.New("db_file not set")
		fmt.Println(err)
		os.Exit(1)
	}
	var err error
	a.db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Printf("open sqlite error %s dbFile %s\n", err.Error(), dbFile)
		os.Exit(1)
	}
	_, err = clMigrate.Migrate(a.db.DB)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
	default:
		if err != nil {
			fmt.Printf("Migrate error: %s dbFile %s\n", err.Error(), dbFile)
			os.Exit(1)
		}
	}
	storePath, err := cfg.UsrCfgDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := storage.NewStorage(a.db, storePath)
	a.srv = service.NewService(r)

	a.addRootCmd().
		addConfigCmd().
		addSaveCmd().
		addDeleteCmd().
		addListCmd().
		addProfileCmd()
	return
}

func (a *app) Execute() {
	defer func() {
		err := a.db.Close()
		if err != nil {
			fmt.Printf("close db error: %s", err)
		}
	}()
	err := a.root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.Glob.GetBool("autosave") && cfg.Glob.Get("changed_at") != nil {
		fmt.Print("Saving global cfg files at exit..")
		err = cfg.Glob.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(" ..Success")
	}
	if cfg.User.Get("name") != nil && cfg.User.GetBool("autosave") && cfg.User.Get("changed_at") != nil {
		fmt.Print("Saving user cfg files at exit..")
		err = cfg.User.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(" ..Success")
	}
}

func GenerateFlags(in interface{}, fs *pflag.FlagSet) error {
	// thanks https://stackoverflow.com/questions/72891199/procedurally-bind-struct-fields-to-command-line-flag-values-using-reflect
	rv := reflect.ValueOf(in)
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
