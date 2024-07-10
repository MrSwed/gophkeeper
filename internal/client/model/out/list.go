package out

import "gophKeeper/internal/client/storage"

type List struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

func (l *List) FromDBItems(r ...storage.DBItem) {
	l.Items = make([]Item, len(r))
	for idx, item := range r {
		l.Items[idx].DBItem = item
	}
}
