package storage

import (
	_ "github.com/mattn/go-sqlite3"
)

type DBItem struct {
	Key         string `db:"key"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type DBRecord struct {
	DBItem
	Filename string `db:"filename"`
}

type StoreRecord struct {
	DBItem
	Data []byte
}

type ListItem struct {
	DBItem
}
