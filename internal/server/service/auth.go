package service

import (
	"context"
	"time"

	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/constant"
	errs "gophKeeper/internal/server/errors"
	"gophKeeper/internal/server/model"
	"gophKeeper/internal/server/repository"

	"golang.org/x/crypto/bcrypt"
)

var _ Auth = (*auth)(nil)

type auth serv

func NewServiceAuth(r repository.Storage, c *config.Config) Auth {
	return &auth{
		r: r,
		c: c,
	}
}

// GetClientToken
// auth by email and password, return client token for sync
func (s auth) GetClientToken(ctx context.Context, req model.AuthRequest) (token []byte, err error) {
	var (
		email    = req.Email
		passHash []byte
		u        model.DBUser
		exp      = time.Now().Add(constant.ExpDuration)
	)
	u, err = s.r.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}
	passHash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	if nil != bcrypt.CompareHashAndPassword(passHash, u.Password) {
		err = errs.ErrorWrongAuth
		return
	}

	token, err = s.r.NewUserClientToken(ctx, u.ID, &exp, req.Meta)
	return
}
