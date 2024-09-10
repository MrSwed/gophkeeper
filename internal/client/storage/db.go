package storage

import (
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
	if query.Limit != 0 {
		b = b.Limit(query.Limit)
	}
	if query.Offset != 0 {
		b = b.Offset(query.Offset)
	}
	return b
}

func (s *dbStore) List(query model.ListQuery) (data []DBItem, err error) {
	var (
		builder = sq.Select("key", "description", "created_at", "updated_at").
			From("storage")
		sql  string
		args []interface{}
	)

	sql, args, err = s.querySqlBuilder(builder, query).ToSql()
	if err != nil {
		return
	}
	err = s.db.Select(&data, sql, args...)
	return
}

func (s *dbStore) Count(query model.ListQuery) (n int, err error) {
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

func (s *dbStore) Get(key string) (DBRecord, error) {
	var data DBRecord
	err := s.db.Get(&data,
		`SELECT key, description, created_at, updated_at, filename, blob FROM storage where key = ?`,
		key)
	if err != nil {
		return DBRecord{}, err
	}
	return data, nil
}
func (s *dbStore) Save(data DBRecord) (err error) {
	_, err = s.db.Exec(`insert into storage 
 (key, description, created_at, updated_at, filename, blob)
 values(?,?,DATETIME('now'),DATETIME('now'),?,?)
 on conflict (key) do update 
  set description=excluded.description,
      updated_at=excluded.updated_at,
      filename=excluded.filename,
      blob=excluded.blob`,
		data.Key, data.Description, data.Filename, data.Blob)
	return
}

func (s *dbStore) Delete(key string) (err error) {
	_, err = s.db.Exec(`delete from storage where key = ?`, key)
	return
}
