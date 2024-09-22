package service

import (
	"context"
	"gophKeeper/internal/helper"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ User = (*user)(nil)

type user serv

func NewServiceUser(r repository.Storage, c *config.Config) User {
	return &user{
		r: r,
		c: c,
	}
}

func (s *user) GetSelf(ctx context.Context) (user model.User, err error) {
	var u model.DBUser
	u, err = s.r.GetUserSelf(ctx)
	if err != nil {
		return
	}
	user.Email = u.Email
	user.PackedKey = u.PackedKey
	user.Description = u.Description
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt
	user.ID = u.ID
	return
}

func (s *user) SaveSelf(ctx context.Context, user *model.User) (err error) {
	u := model.DBUser{
		ID:          user.ID,
		Email:       user.Email,
		Description: user.Description,
		PackedKey:   user.PackedKey,
	}
	if user.Password != "" {
		u.Password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
	}

	err = s.r.SaveUser(ctx, &u)
	if err == nil {
		user.ID = u.ID
		user.CreatedAt = u.CreatedAt
		user.UpdatedAt = u.UpdatedAt
	}
	return
}

func (s *user) UserIDByToken(ctx context.Context, token []byte) (id uuid.UUID, err error) {
	return s.r.GetUserIDByToken(ctx, token)
}

func (s *user) DeleteSelf(ctx context.Context) (err error) {
	var userID uuid.UUID
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}

	return s.r.DeleteUser(ctx, userID)
}
