package model

import (
	"time"
)

type ItemSync struct {
	DBItem
	Blob []byte `db:"blob" json:"blob"`
}

type ListRequest struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type ListResponse struct {
	Items []DBItem `json:"items"`
	Total int64    `json:"total"`
}

type OkResponse struct {
	Ok bool `json:"ok"`
}

type RegisterClientRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClientToken struct {
	Token []byte `json:"token"`
}

type UserSync struct {
	Password    string     `json:"password"`
	Email       string     `json:"email"`
	Description *string    `json:"description"`
	PackedKey   []byte     `json:"packedKey"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
