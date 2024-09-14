package repository

import (
	"context"
	"database/sql"
	"errors"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type authStore store

var _ UserStorage = (*authStore)(nil)

const (
	userTableName   = "users"
	clientTableName = "clients"
)

func NewAuthStorage(c *config.StorageConfig, db *sqlx.DB) *authStore {
	return &authStore{
		db: db,
		c:  c,
	}
}

func (s *authStore) NewUser(ctx context.Context, user model.User) (userID uuid.UUID, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Insert(clientTableName).
		SetMap(map[string]any{
			"email":       user.Email,
			"password":    user.Password,
			"description": user.Description,
			"packed_key":  user.PackedKey,
		}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &userID, query, args...)
	return
}

func (s *authStore) GetUserByID(ctx context.Context, userID uuid.UUID) (user model.User, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`id, description, email, packed_key, created_at, updated_at`).
		From(userTableName).
		Where("id = ?", userID).ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &user, query, args...)
	return
}

func (s *authStore) GetUserIDByToken(ctx context.Context, token []byte) (userID uuid.UUID, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`u.id`).
		LeftJoin(clientTableName+" c on c.user_id = u.id").
		From(userTableName+" u").
		Where("token = ?", token).ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &userID, query, args...)

	return
}

func (s *authStore) GetByEmail(ctx context.Context, email string) (user model.User, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`id, email, description, createdAt, updatedAt`).
		From(userTableName).
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &user, query, args...)
	return
}

func (s *authStore) NewUserClientToken(ctx context.Context, userID uuid.UUID, expAt *time.Time, meta any) (token []byte, err error) {
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
		Suffix("RETURNING token").
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &token, query, args...)
	return
}

func (s *authStore) DeleteClient(ctx context.Context, token []byte) (err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Delete(clientTableName).
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return
	}
	_, err = s.db.ExecContext(ctx, query, args...)
	return
}

func (s *authStore) DeleteUser(ctx context.Context, userID uuid.UUID) (err error) {
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

	return
}
