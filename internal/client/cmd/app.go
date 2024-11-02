package cmd

import (
	"errors"
	"fmt"

	cfg "gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/service"
	"gophKeeper/internal/client/storage"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

type App interface {
	Execute() error
}

type BuildMetadata struct {
	Version string `json:"buildVersion"`
	Date    string `json:"buildDate"`
	Commit  string `json:"buildCommit"`
}

type app struct {
	v    BuildMetadata
	db   *sqlx.DB
	srv  service.Service
	root *cobra.Command
}

func NewApp(b BuildMetadata) (a *app) {
	a = (&app{v: b}).
		addRootCmd().
		addConfigCmd().
		addSaveCmd().
		addViewCmd().
		addDeleteCmd().
		addListCmd().
		addProfileCmd().
		addSyncCmd()
	return
}

// Srv
// initialize the service with the configured store if it is not early
func (a *app) Srv() service.Service {
	if a.srv == nil {
		err := cfg.UserLoad()
		if err != nil {
			return service.NewServiceError(fmt.Errorf("error load config: %v", err))
		}
		// todo: do not print it if loaded early

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

func (a *app) Close() (err error) {
	if a.db != nil {
		defer func() {
			err = a.db.Close()
			if err != nil {
				a.root.Printf("close db error: %s", err)
			}
			a.srv = nil
		}()
	}
	if cfg.Glob.GetBool("autosave") && cfg.Glob.IsChanged() {
		a.root.Print("Saving global cfg files at exit..")
		err = cfg.Glob.Save()
		if err != nil {
			a.root.Println(err)
			// return
		}
		a.root.Println(" ..Success")
	}
	if cfg.User.Viper != nil && cfg.User.Get("name") != nil && cfg.User.GetBool("autosave") && cfg.User.IsChanged() {
		a.root.Print("Saving user cfg files at exit..")
		err = cfg.User.Save()
		if err != nil {
			a.root.Println(err)
			// return
		}
		a.root.Println(" ..Success")
	}
	return
}

func (a *app) Execute() (err error) {
	defer func() {
		err = errors.Join(err, a.Close())
	}()

	err = a.root.Execute()
	if err != nil {
		a.root.Println(err)
		return
	}

	return
}
