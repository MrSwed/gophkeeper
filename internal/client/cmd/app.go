package cmd

import (
	"errors"
	"fmt"
	cfg "gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/service"
	"gophKeeper/internal/client/storage"
	"gophKeeper/internal/helper"
	"os"
	"reflect"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

type App interface {
	Execute() error
}

type app struct {
	db    *sqlx.DB
	srv   service.Service
	root  *cobra.Command
	debug bool
}

func NewApp() (a *app) {
	var err error
	a = &app{}
	err = cfg.UserLoad()
	if err != nil {
		fmt.Printf("Error load current user profile: %s", err)
		os.Exit(1)
	}

	dbFile := cfg.User.GetString("db_file")
	if dbFile == "" {
		err := errors.New("db_file not set")
		fmt.Println(err)
		os.Exit(1)
	}
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

func GenFlags(in interface{}) (flags []string, err error) {
	rv := reflect.ValueOf(in)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		err = errors.New("not pointer-to-a-struct") // exit if not pointer-to-a-struct
		return
	}
	rv = rv.Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		tagNames := [2]string{}
		copy(tagNames[:], strings.SplitN(sf.Tag.Get(("flag")), ",", 2))
		flags = append(flags, tagNames[0])
	}
	return
}

func modelGenerateFlags(dst any, cmd *cobra.Command, debug *bool) (err error) {
	if debug != nil {
		cmd.Flags().BoolVarP(debug, "debug", "", *debug, "debug flag")
	}
	if common, ok := dst.(model.Base); ok {
		err = helper.GenerateFlags(common.GetBase(), cmd.Flags())
	}
	if data, ok := dst.(model.Data); ok {
		err = errors.Join(err, helper.GenerateFlags(data.GetDst(), cmd.Flags()))
	}
	return
}
