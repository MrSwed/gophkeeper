package cmd

import (
	"errors"
	"fmt"
	cfg "gophKeeper/internal/client/config"
	clMigrate "gophKeeper/internal/client/migrate"
	"gophKeeper/internal/client/service"
	"gophKeeper/internal/client/storage"
	"os"

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
		addListCmd()
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
