package storage

import (
	"github.com/jmoiron/sqlx"
)

type DB interface {
	List() (data []ListItem, err error)
	Get(key string) (data DBRecord, err error)
	Set(data DBRecord) (err error)
	Delete(key string) (err error)
}

type File interface {
	Get(fileName string) (b []byte, err error)
	Set(fileName string, b []byte) (err error)
	Delete(fileName string) (err error)
}

type Storage struct {
	DB   DB
	File File
}

func NewStorage(db *sqlx.DB, path string) *Storage {
	return &Storage{
		DB:   NewDBStore(db),
		File: NewFileStore(path),
	}
}
