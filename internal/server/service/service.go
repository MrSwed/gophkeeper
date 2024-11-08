package service

import (
	"context"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/repository"

	"github.com/google/uuid"
)

type serv struct {
	r repository.Storage
	c *config.Config
}

type Data interface {
	ListSelf(ctx context.Context, q *model.ListQuery) (list model.List, err error)
	GetSelfItem(ctx context.Context, k string) (item *model.Item, err error)
	SaveSelfItem(ctx context.Context, item *model.Item) (err error)
}

type Auth interface {
	GetClientToken(ctx context.Context, req model.AuthRequest) (token []byte, err error)
}

type User interface {
	GetSelf(ctx context.Context) (user model.User, err error)
	SaveSelf(ctx context.Context, user *model.User) (err error)
	DeleteSelf(ctx context.Context) (err error)

	UserIDByToken(ctx context.Context, token []byte) (uuid.UUID, error)
}

type Service interface {
	Auth
	Data
	User
}

type service struct {
	Auth
	Data
	User
}

func New(r repository.Storage, c *config.Config) Service {
	return &service{
		Data: NewServiceData(r, c),
		Auth: NewServiceAuth(r, c),
		User: NewServiceUser(r, c),
	}
}
