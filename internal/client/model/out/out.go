package out

import (
	"encoding/json"
	"gophKeeper/internal/client/model"
)

type dTypeRaw []byte

func (d *dTypeRaw) UnmarshalJSON(b []byte) (err error) {
	*d = b
	return
}

type dRaw struct {
	Type string   `json:"type"`
	Data dTypeRaw `json:"data"`
}

type Item struct {
	model.DBItem
	Data model.Data `json:"data"`
}

func (i *Item) UnmarshalJSON(b []byte) (err error) {
	t := new(dRaw)
	if err = json.Unmarshal(b, &t); err != nil {
		return
	}
	i.Data, err = model.GetNewDataModel(t.Type)
	if err != nil {
		return
	}
	err = json.Unmarshal(t.Data, &i.Data)

	return
}

func (i *Item) FromDBItem(dbItem model.DBItem) {
	i.DBItem = dbItem
}

type List struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

func (l *List) FromDBItems(r ...model.DBItem) {
	l.Items = make([]Item, len(r))
	for idx, item := range r {
		l.Items[idx].FromDBItem(item)
	}
}
