package model

import (
	"time"
)

type ItemShort struct {
	Key         string     `db:"key" json:"key"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type Item struct {
	ItemShort
	Blob []byte `db:"blob" json:"blob"`
}

type List struct {
	Items []ItemShort `json:"items"`
	Total int64       `json:"total"`
}

type DBRecord struct {
	ItemShort
	FileName *string `db:"filename,omitempty"`
	Blob     []byte  `db:"blob"`
}

func (i *Item) IsNew() bool {
	return i.CreatedAt.IsZero()
}
