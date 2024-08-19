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
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
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

	a.addRootCmd().
		addConfigCmd().
		addSaveCmd().
		addViewCmd().
		addDeleteCmd().
		addListCmd().
		addProfileCmd()
	return
}

// Srv
// initialize the service with the configured store if it is not early
func (a *app) Srv() service.Service {
	if a.srv == nil {
		err := cfg.GlobalLoad()
		if err != nil {
			a.root.Println(err)
		}
		err = cfg.UserLoad()
		if err != nil {
			return service.NewServiceError(fmt.Errorf("error load current user profile: %v", err))
		}
		a.root.Printf("User %s configuration loaded\n", cfg.User.GetString("name"))

		dbFile := cfg.User.GetString("db_file")
		if dbFile == "" {
			return service.NewServiceError(errors.New("error db_file - is not set"))
		}
		a.db, err = sqlx.Open("sqlite3", dbFile)
		if err != nil {
			return service.NewServiceError(fmt.Errorf("open sqlite error %s dbFile %s\n", err.Error(), dbFile))
		}
		_, err = clMigrate.Migrate(a.db.DB)
		switch {
		case errors.Is(err, migrate.ErrNoChange):
		default:
			if err != nil {
				return service.NewServiceError(fmt.Errorf("db update error: %s dbFile %s\n", err.Error(), dbFile))
			}
		}
		storePath, err := cfg.UsrCfgDir()
		if err != nil {
			return service.NewServiceError(fmt.Errorf("usrCfgDir error: %s \n", err))
		}
		a.srv = service.NewService(storage.NewStorage(a.db, storePath))
	}
	return a.srv
}

func (a *app) Close() {
	if a.db != nil {
		defer func() {
			err := a.db.Close()
			if err != nil {
				a.root.Printf("close db error: %s", err)
			}
		}()
	}
}

func (a *app) Execute() {
	defer a.Close()

	err := a.root.Execute()
	if err != nil {
		a.root.Println(err)
		os.Exit(1)
	}

	if cfg.Glob.GetBool("autosave") && cfg.Glob.Get("changed_at") != nil {
		a.root.Print("Saving global cfg files at exit..")
		err = cfg.Glob.Save()
		if err != nil {
			a.root.Println(err)
			os.Exit(1)
		}
		a.root.Println(" ..Success")
	}
	if cfg.User.Viper != nil && cfg.User.Get("name") != nil && cfg.User.GetBool("autosave") && cfg.User.Get("changed_at") != nil {
		a.root.Print("Saving user cfg files at exit..")
		err = cfg.User.Save()
		if err != nil {
			a.root.Println(err)
			os.Exit(1)
		}
		a.root.Println(" ..Success")
	}
}

func GenFlags(in any) (flags []string, err error) {
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
		copy(tagNames[:], strings.SplitN(sf.Tag.Get("flag"), ",", 2))
		if tagNames[0] != "" {
			flags = append(flags, tagNames[0])
		}
	}
	return
}
