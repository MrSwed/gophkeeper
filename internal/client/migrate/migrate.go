package migrate

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:generate migrate create -ext sql -dir sql -format 20060102030405 gopher_keeper

//go:embed sql/*.sql
var FS embed.FS

func Migrate(db *sql.DB) (version [2]uint, err error) {
	var d source.Driver
	if d, err = iofs.New(FS, "sql"); err != nil {
		return
	}
	defer d.Close()
	var driver database.Driver
	if driver, err = sqlite.WithInstance(db, &sqlite.Config{}); err != nil {
		return
	}
	var m *migrate.Migrate
	if m, err = migrate.NewWithInstance("iofs", d, "", driver); err != nil {
		return
	}
	version[0], _, _ = m.Version()
	err = m.Up()
	version[1], _, _ = m.Version()
	return
}
