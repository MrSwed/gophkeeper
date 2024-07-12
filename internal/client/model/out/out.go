package out

import "gophKeeper/internal/client/storage"

type Item struct {
	storage.DBItem
	Data any `json:"data"`
}

func (i *Item) FromDBItem(dbItem storage.DBItem) {
	i.DBItem = dbItem
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Bin struct {
	Bin []byte `json:"bin"`
}

type Card struct {
	Exp    string `json:"exp"`
	Number string `json:"number"`
	Name   string `json:"name"`
	CVV    string `json:"cvv"`
}

type Text struct {
	Text string `json:"text"`
}
