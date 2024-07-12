package out

import (
	"encoding/json"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/storage"
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
	storage.DBItem
	Data any `json:"data"`
}

func (i *Item) UnmarshalJSON(b []byte) (err error) {
	t := new(dRaw)
	if err = json.Unmarshal(b, &t); err != nil {
		return
	}
	i.Data, err = model.GetNewModel(t.Type)
	if err != nil {
		return
	}
	err = json.Unmarshal(t.Data, &i.Data)

	return
}

func (i *Item) FromDBItem(dbItem storage.DBItem) {
	i.DBItem = dbItem
}