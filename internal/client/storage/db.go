package storage

import (
	"time"

	"gophKeeper/internal/client/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	sq "github.com/Masterminds/squirrel"
)

type dbStore struct {
	db *sqlx.DB
}

func NewDBStore(db *sqlx.DB) *dbStore {
	return &dbStore{
		db: db,
	}
}

func (s *dbStore) querySqlBuilder(b sq.SelectBuilder, query model.ListQuery) sq.SelectBuilder {
	if query.Key != "" {
		b = b.Where(sq.Like{"key": "%" + query.Key + "%"})
	}
	if query.Description != "" {
		b = b.Where(sq.Like{"description": "%" + query.Description + "%"})
	}
	if query.CreatedAt != "" {
		b = b.Where(sq.Like{"created_at": "%" + query.CreatedAt + "%"})
	}
	if query.UpdatedAt != "" {
		b = b.Where(sq.Like{"updated_at": "%" + query.UpdatedAt + "%"})
	}
	if query.SyncAt != "" {
		b = b.Where("sync_at is null or sync_at < ?", query.SyncAt)
	}
	if !query.Deleted {
		b = b.Where("blob is not null or filename is not null")
	}
	return b
}

func (s *dbStore) List(query model.ListQuery) (data []model.DBItem, err error) {
	var (
		builder = sq.Select("key", "description", "created_at", "updated_at", "sync_at").
			From("storage")
		sql  string
		args []interface{}
	)
	if query.Limit != 0 {
		builder = builder.Limit(query.Limit)
	}
	if query.Offset != 0 {
		builder = builder.Offset(query.Offset)
	}
	if query.OrderBy != "" {
		builder = builder.OrderBy(query.OrderBy)
	}

	sql, args, err = s.querySqlBuilder(builder, query).ToSql()
	if err != nil {
		return
	}
	err = s.db.Select(&data, sql, args...)
	return
}

func (s *dbStore) Count(query model.ListQuery) (n uint64, err error) {
	var (
		builder = sq.Select("count(*) as count").
			From("storage")
		sql  string
		args []interface{}
	)
	query.Limit = 0
	query.Offset = 0
	sql, args, err = s.querySqlBuilder(builder, query).ToSql()
	if err != nil {
		return
	}
	err = s.db.Get(&n, sql, args...)

	return
}

func (s *dbStore) Get(key string) (model.DBRecord, error) {
	var data model.DBRecord
	err := s.db.Get(&data,
		`SELECT key, description, created_at, updated_at, filename, blob, sync_at FROM storage where key = ?`,
		key)
	if err != nil {
		return model.DBRecord{}, err
	}
	return data, nil
}

func (s *dbStore) Save(data model.DBRecord) (err error) {
	createdAt := data.CreatedAt.Format(time.DateTime)
	if data.CreatedAt.IsZero() {
		createdAt = time.Now().Format(time.DateTime)
	}
	var updatedAt *string
	if data.UpdatedAt != nil {
		updatedAt = &[]string{data.UpdatedAt.Format(time.DateTime)}[0]
	}
	_, err = s.db.Exec(`insert into storage 
 (key, description, created_at, updated_at, filename, blob, sync_at)
 values(?,?,?,?,?,?,?)
 on conflict (key) do update 
  set description=excluded.description,
      updated_at=case when excluded.updated_at is not null then excluded.updated_at else DATETIME('now','localtime') end,
      filename=excluded.filename,
      blob=excluded.blob,
      sync_at=excluded.sync_at`,
		data.Key, data.Description, createdAt, updatedAt, data.Filename, data.Blob, data.SyncAt)
	return
}

func (s *dbStore) Delete(key string) (err error) {
	_, err = s.db.Exec(`update storage 
set description = '', updated_at=DATETIME('now','localtime'), filename = null, blob=null
where key = ?`, key)
	return
}
