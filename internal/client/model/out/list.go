package out

import "gophKeeper/internal/client/storage"

type List struct {
	Items []storage.ListItem `json:"items"`
	Total int                `json:"total"`
}

type Item struct {
	storage.DBItem
	Data []byte `json:"data"`
}
