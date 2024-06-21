package model

import "time"

type Common struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	createdDate time.Time
	updatedDate time.Time
}

type Auth struct {
	Common
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Text struct {
	Common
	Text string `json:"text"`
}

type Bin struct {
	Common
	Bin []byte `json:"bin"`
}

type Card struct {
	Year  uint8 `json:"year"`
	Month uint8 `json:"month"`
	Common
	Number string `json:"number"`
	Name   string `json:"name"`
	CVV    string `json:"cvv"`
}
