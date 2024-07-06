package storage

import (
	_ "github.com/mattn/go-sqlite3"
)

type DBItem struct {
	Key         string  `db:"key" json:"key"`
	Description *string `db:"description" json:"description"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	UpdatedAt   string  `db:"updated_at" json:"updated_at"`
}

type DBRecord struct {
	DBItem
	Filename string `db:"filename"`
}

type ListItem struct {
	DBItem
}
