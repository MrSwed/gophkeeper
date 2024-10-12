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

type List struct {
	Items []model.DBItem `json:"items"`
	Total uint64         `json:"total"`
}
