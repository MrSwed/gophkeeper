package model

import (
	"time"

	"github.com/google/uuid"
)

type PassChangeRequest struct {
	Password string `json:"password" validate:"omitempty,password"`
}

func (p *PassChangeRequest) Validate(fields ...string) error {
	return ValidateStruct(p, fields...)
}

type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	Description *string
	PackedKey   []byte
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type DBUser struct {
	ID          uuid.UUID  `db:"id"`
	Email       string     `db:"email"`
	Password    []byte     `db:"password"`
	Description *string    `db:"description"`
	PackedKey   []byte     `db:"packed_key"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
