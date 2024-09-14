package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Password    string
	Email       string
	Description string
	PackedKey   []byte
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type DBUser struct {
	ID          uuid.UUID  `db:"id"`
	Password    []byte     `db:"password"`
	Email       string     `db:"email"`
	Description string     `db:"description"`
	PackedKey   []byte     `db:"packed_key"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
