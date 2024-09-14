package model

import (
	"time"

	"github.com/google/uuid"
)

type ItemShort struct {
	Key         string     `db:"key"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
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
	UserID   uuid.UUID `db:"user_id" validate:"required"`
	FileName *string   `db:"filename,omitempty"`
	Blob     []byte    `db:"blob"`
}

func (i *Item) IsNew() bool {
	return i.CreatedAt.IsZero()
}
