package storage

import (
	"gophKeeper/internal/client/model"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	List(query model.ListQuery) (data []DBItem, err error)
	Count(query model.ListQuery) (n int, err error)
	Get(key string) (data DBRecord, err error)
	Set(data DBRecord) (err error)
	Delete(key string) (err error)
}

type File interface {
	Get(fileName string) (b []byte, err error)
	Save(fileName string, b []byte) (err error)
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
