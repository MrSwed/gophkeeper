package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Password    []byte     `json:"password" db:"password"`
	Email       string     `json:"email" db:"email"`
	Description string     `json:"description" db:"description"`
	PackedKey   []byte     `json:"packedKey" db:"packed_key"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
}
