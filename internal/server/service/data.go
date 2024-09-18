package service

import (
	"context"
	"gophKeeper/internal/helper"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/repository"

	"github.com/google/uuid"
)

var _ Data = (*serv)(nil)

// NewServiceData return main service methods
func NewServiceData(r repository.Storage, c *config.Config) *serv {
	return &serv{r: r, c: c}
}

func (s *serv) ListSelf(ctx context.Context, q *model.ListQuery) (list model.List, err error) {
	if err = q.Validate(); err != nil {
		return
	}
	if q == nil {
		q = &model.ListQuery{}
	}
	q.UserID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}

	if list.Total, err = s.r.CountDataItems(ctx, q); err != nil {
		return
	}
	list.Items, err = s.r.ListDataItems(ctx, q)
	return
}

func (s *serv) GetSelfItem(ctx context.Context, k string) (item *model.Item, err error) {
	item = &model.Item{ItemShort: model.ItemShort{Key: k}}
	var (
		userID uuid.UUID
		dbItem model.DBRecord
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}

	dbItem, err = s.r.GetDataItem(ctx, userID, k)
	if err != nil {
		return
	}
	item.ItemShort = dbItem.ItemShort
	item.Blob = dbItem.Blob
	// if dbItem.FileName != nil && item.Blob == nil {
	// item.Blob, err = s.r.GetFile(*dbItem.FileName)
	// }

	return
}

func (s *serv) SaveSelfItem(ctx context.Context, item *model.Item) (err error) {
	var (
		dbItem model.DBRecord
	)
	dbItem.UserID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}

	dbItem.ItemShort = item.ItemShort
	// todo: save to file if size exceeds max.
	// dbItem.Blob = nil
	// dbItem.FileName = getFileStorageName
	dbItem.Blob = item.Blob

	err = s.r.SaveDataItem(ctx, dbItem)

	return
}
