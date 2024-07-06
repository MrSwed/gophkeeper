package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type dbStore struct {
	db *sqlx.DB
}

func NewDBStore(db *sqlx.DB) *dbStore {
	return &dbStore{
		db: db,
	}
}

func (s *dbStore) List() (data []ListItem, err error) {
	err = s.db.Select(&data,
		`SELECT key, description, created_at, updated_at FROM storage`)
	return
}
func (s *dbStore) Get(key string) (data DBRecord, err error) {
	err = s.db.Get(&data,
		`SELECT key, description, created_at, updated_at, filename FROM storage where key = ?`,
		key)
	return
}
func (s *dbStore) Set(data DBRecord) (err error) {
	_, err = s.db.Exec(`insert into storage 
 (key, description, created_at, updated_at, filename) values(?,?,?,?,?)
 on conflict (key) do update 
  set description=excluded.description,
      updated_at=excluded.updated_at,
      filename=excluded.filename`,
		data.Key, data.Description, time.Now(), time.Now(), data.Filename)
	return
}

func (s *dbStore) Delete(key string) (err error) {
	_, err = s.db.Exec(`delete from storage where key = ?`, key)
	return
}
