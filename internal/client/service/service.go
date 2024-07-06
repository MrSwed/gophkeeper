package service

import (
	"database/sql"
	"errors"
	"gophKeeper/internal/client/model/input"
	"gophKeeper/internal/client/model/out"
	"gophKeeper/internal/client/storage"
	"time"
)

type Service interface {
	List(query input.ListQuery) (data out.List, err error)
	Get(key string) (data out.Item, err error)
	Set(data input.Model) (err error)
	Delete(key string) (err error)
}

type service struct {
	r *storage.Storage
}

func NewService(r *storage.Storage) *service {
	return &service{r: r}
}

func (s *service) List(query input.ListQuery) (data out.List, err error) {
	if err = query.Validate(); err != nil {
		return
	}
	if data.Total, err = s.r.DB.Count(query); err != nil {
		return
	}
	if data.Items, err = s.r.DB.List(query); err != nil {
		return
	}
	return
}

func (s *service) Get(key string) (data out.Item, err error) {
	var r storage.DBRecord
	if r, err = s.r.DB.Get(key); err != nil {
		return
	}
	if data.Data, err = s.r.File.Get(r.Filename); err != nil {
		return
	}
	data.DBItem = r.DBItem

	return
}

func (s *service) Set(data input.Model) (err error) {
	if err = data.Validate(); err != nil {
		return
	}
	var r storage.DBRecord
	if r, err = s.r.DB.Get(data.GetKey()); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return
	}
	if r.Key == "" {
		r.Key = data.GetKey()
		r.Filename = time.Now().Format("20060102150405") + "-" + r.Key
	}
	r.Description = data.GetDescription()
	var b []byte

	if b, err = data.Bytes(); err != nil {
		return
	}

	// todo: crypt here ??

	if err = s.r.File.Save(r.Filename, b); err != nil {
		return
	}

	err = s.r.DB.Set(r)
	return
}

func (s *service) Delete(key string) (err error) {
	var r storage.DBRecord
	if r, err = s.r.DB.Get(key); err != nil {
		return
	}
	if err = s.r.File.Delete(r.Filename); err != nil {
		return
	}
	err = s.r.DB.Delete(key)
	return
}
