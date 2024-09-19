package repository

import (
	"context"
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

type dataStore store

func (s *dataStore) GetDataItem(ctx context.Context, userID uuid.UUID, key string) (item model.DBRecord, err error) {
	var (
		query string
		args  []interface{}
	)
	query, args, err = sq.Select(`key, description, created_at, updated_at, filename, blob`).
		From(storeTableName).
		Where("key = ?", key).
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &item, query, args...)
	return
}

func (s *dataStore) ListDataItems(ctx context.Context, q *model.ListQuery) (list []model.ItemShort, err error) {
	var (
		query string
		args  []interface{}
	)

	sqlBuild := sq.Select("key", "description", "created_at", "updated_at").
		From(storeTableName)
	if q != nil {
		if q.UserID != uuid.Nil {
			sqlBuild = sqlBuild.Where("user_id = ?", q.UserID)
		}
		if q.Limit != 0 {
			sqlBuild = sqlBuild.Limit(q.Limit)
		}
		if q.Offset != 0 {
			sqlBuild = sqlBuild.Offset(q.Offset)
		}
	}
	query, args, err = sqlBuild.ToSql()
	if err != nil {
		return
	}
	err = s.db.SelectContext(ctx, &list, query, args...)

	return
}

func (s *dataStore) CountDataItems(ctx context.Context, q *model.ListQuery) (total int64, err error) {
	var (
		query string
		args  []interface{}
	)
	sqlBuild := sq.Select("count(*)").From(storeTableName)
	if q != nil {
		if q.UserID != uuid.Nil {
			sqlBuild = sqlBuild.Where("user_id = ?", q.UserID)
		}
	}
	query, args, err = sqlBuild.ToSql()
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &total, query, args...)

	return
}

func (s *dataStore) SaveDataItem(ctx context.Context, item model.DBRecord) (err error) {
	var (
		query string
		args  []interface{}
	)
	// todo: use it if db cannot check key reference
	// if item.UserID == uuid.Nil {
	// 	err = errors.New("UserID should not be nil")
	// 	return
	// }

	query, args, err = sq.Insert(storeTableName).
		Columns(`key, user_id, description, created_at, updated_at, filename, blob`).
		Values(item.Key, item.UserID, item.Description, item.CreatedAt, item.UpdatedAt, item.FileName, item.Blob).
		Suffix(`on conflict (key, user_id) do update 
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
