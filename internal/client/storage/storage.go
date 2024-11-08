package storage

import (
	"gophKeeper/internal/client/model"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	List(query model.ListQuery) (data []model.DBItem, err error)
	Count(query model.ListQuery) (n uint64, err error)
	Get(key string) (data model.DBRecord, err error)
	Save(data model.DBRecord) (err error)
	Delete(key string) (err error)
}

type File interface {
	GetStored(fileName string) (b []byte, err error)
	SaveStore(fileName string, b []byte) (err error)
	Delete(fileName string) (err error)
	GetOrigin(filePath string) (b []byte, err error)
	SaveOrigin(filePath string, b []byte) (err error)
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
