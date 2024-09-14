package repository

import (
	"context"
	"gophKeeper/internal/helper"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewDBRepository(c *config.StorageConfig, db *sqlx.DB) *dataStore {
	return &dataStore{
		db: db,
		c:  c,
	}
}

var _ DataStorage = (*dataStore)(nil)

const storeTableName = "storage"

type dataStore store

func (s *dataStore) GetItem(ctx context.Context, key string) (item *model.DBRecord, err error) {
	var (
		userID uuid.UUID
		query  string
		args   []interface{}
	)
	item = new(model.DBRecord)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}
	query, args, err = sq.Select(`key, description, created_at, updated_at, filename, blob`).
		From(storeTableName).
		Where("key = ?", key).
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.Get(&item, query, args...)
	return
}

func (s *dataStore) ListItems(ctx context.Context, q *model.ListQuery) (list []model.ItemShort, err error) {
	var (
		userID uuid.UUID
		query  string
		args   []interface{}
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}
	sqlBuild := sq.Select("key", "description", "created_at", "updated_at").
		From(storeTableName).
		Where("user_id = ?", userID)
	if q.Limit != 0 {
		sqlBuild = sqlBuild.Limit(q.Limit)
	}
	if q.Offset != 0 {
		sqlBuild = sqlBuild.Offset(q.Offset)
	}

	query, args, err = sqlBuild.
		ToSql()
	if err != nil {
		return
	}
	err = s.db.Select(&list, query, args...)
	return
}

func (s *dataStore) CountItems(ctx context.Context, _ *model.ListQuery) (total int, err error) {
	var (
		userID uuid.UUID
		query  string
		args   []interface{}
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}
	sqlBuild := sq.Select("count(*)").
		From(storeTableName).
		Where("user_id = ?", userID)

	query, args, err = sqlBuild.
		ToSql()
	if err != nil {
		return
	}
	err = s.db.Get(&total, query, args...)
	return

}

func (s *dataStore) SaveItem(ctx context.Context, item *model.DBRecord) (err error) {
	var (
		userID uuid.UUID
		query  string
		args   []interface{}
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}
	query, args, err = sq.Insert(storeTableName).
		Columns(`key, user_id, description, created_at, updated_at, filename, blob`).
		Values(item.Key, userID, item.Description, item.CreatedAt, item.UpdatedAt, item.FileName, item.Blob).
		Prefix(`on conflict (key, userID) do update 
  set description=excluded.description,
      updated_at=excluded.updated_at,
      filename=excluded.filename,
      blob=excluded.blob`).
		ToSql()
	if err != nil {
		return
	}
	_, err = s.db.ExecContext(ctx, query, args...)
	return
}

func (s *dataStore) DeleteItem(ctx context.Context, key string) (err error) {
	var (
		userID uuid.UUID
		query  string
		args   []interface{}
	)
	userID, err = helper.GetCtxUserID(ctx)
	if err != nil {
		return
	}
	query, args, err = sq.Delete(storeTableName).
		Where("user_id = ?", userID).
		Where("key = ?", key).
		ToSql()
	if err != nil {
		return
	}
	_, err = s.db.ExecContext(ctx, query, args...)
	return
}
