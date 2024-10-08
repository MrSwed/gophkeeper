package service

import (
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/out"
)

var _ Service = (*serviceError)(nil)

type serviceError struct {
	e error
}

func NewServiceError(e error) *serviceError {
	return &serviceError{e: e}
}

func (s *serviceError) ChangePasswd() (err error) {
	err = s.e
	return
}

func (s *serviceError) GetToken() (token string, err error) {
	err = s.e
	return
}

func (s *serviceError) List(_ model.ListQuery) (data out.List, err error) {
	err = s.e
	return
}

func (s *serviceError) Get(_ string) (data out.Item, err error) {
	err = s.e
	return
}

func (s *serviceError) GetRaw(_ string) (data model.DBRecord, err error) {
	err = s.e
	return
}

func (s *serviceError) Save(_ model.Model) (err error) {
	err = s.e
	return
}

func (s *serviceError) SaveRaw(_ model.DBRecord) (err error) {
	err = s.e
	return
}

func (s *serviceError) Delete(_ string) (err error) {
	err = s.e
	return
}
