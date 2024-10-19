package repository

import (
	"context"
	"database/sql"
	"errors"
	"gophKeeper/internal/helper"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userStore store

var _ UserStorage = (*userStore)(nil)

func NewUserStorage(c *config.StorageConfig, db *sqlx.DB) *userStore {
	return &userStore{
		db: db,
		c:  c,
	}
}

func (s *userStore) SaveUser(ctx context.Context, user *model.DBUser) (err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Insert(userTableName).
		SetMap(map[string]any{
			"email":       user.Email,
			"password":    user.Password,
			"description": user.Description,
			"packed_key":  user.PackedKey,
		}).
		Suffix(`
on conflict (email) do update
set description=excluded.description,
      password=case when excluded.password <> '' then excluded.password else password end,
      packed_key=excluded.packed_key
RETURNING id, created_at, updated_at`).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, user, query, args...)
	return
}

func (s *userStore) GetUserSelf(ctx context.Context) (user model.DBUser, err error) {
	var (
		query  string
		args   []interface{}
		userID uuid.UUID
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}

	query, args, err = sq.Select(`id, description, email, packed_key, created_at, updated_at`).
		From(userTableName).
		Where("id = ?", userID).ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &user, query, args...)
	return
}

func (s *userStore) GetUserIDByToken(ctx context.Context, token []byte) (userID uuid.UUID, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`u.id`).
		LeftJoin(clientTableName+" c on c.user_id = u.id").
		From(userTableName+" u").
		Where("token = decode(?, 'hex');", token).ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &userID, query, args...)

	return
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (user model.DBUser, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`id, email, password, description, created_at, updated_at, packed_key`).
		From(userTableName).
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &user, query, args...)
	return
}

func (s *userStore) NewUserClientToken(ctx context.Context, userID uuid.UUID, expAt *time.Time, meta any) (token []byte, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Insert(clientTableName).
		SetMap(map[string]any{
			"user_id":    userID,
			"expired_at": expAt,
			"meta":       meta,
		}).
		Suffix("RETURNING encode(token,'hex')").
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &token, query, args...)
	return
}

func (s *userStore) DeleteUser(ctx context.Context, userID uuid.UUID) (err error) {
	var (
		query string
		args  []interface{}
		tx    *sqlx.Tx
	)
	tx, err = s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		rErr := tx.Rollback()
		if rErr != nil && !errors.Is(rErr, sql.ErrTxDone) {
			err = errors.Join(err, rErr)
		}
	}()

	query, args, err = sq.Delete(storeTableName).
		Where("user_id = ?", userID).
		ToSql()
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}

	query, args, err = sq.Delete(clientTableName).
		Where("user_id = ?", userID).
		ToSql()
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}

	query, args, err = sq.Delete(userTableName).
		Where("id = ?", userID).
		ToSql()
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
