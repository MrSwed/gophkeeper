package storage

import (
	"gophKeeper/internal/client/model/input"

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

func (s *dbStore) querySqlBuilder(b sq.SelectBuilder, query input.ListQuery) sq.SelectBuilder {
	if query.Key != "" {
		b = b.Where(sq.ILike{"key": query.Key})
	}
	if query.Description != "" {
		b = b.Where(sq.ILike{"description": query.Description})
	}
	if query.CreatedAt != "" {
		b = b.Where(sq.ILike{"created_at": query.CreatedAt})
	}
	if query.UpdatedAt != "" {
		b = b.Where(sq.ILike{"updated_at": query.UpdatedAt})
	}
	if query.Limit != 0 {
		b.Limit(query.Limit)
	}
	if query.Offset != 0 {
		b = b.Offset(query.Offset)
	}
	return b
}

func (s *dbStore) List(query input.ListQuery) (data []ListItem, err error) {
	var (
		builder = sq.Select("key", "description", "created_at", "updated_at")
		sql     string
		args    []interface{}
	)

	sql, args, err = s.querySqlBuilder(builder, query).ToSql()
	if err != nil {
		return
	}
	err = s.db.Select(&data, sql, args...)
	return
}

func (s *dbStore) Count(query input.ListQuery) (n int, err error) {
	var (
		builder = sq.Select("count(*) as count")
		sql     string
		args    []interface{}
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

func (s *dbStore) Get(key string) (data DBRecord, err error) {
	err = s.db.Get(&data,
		`SELECT key, description, created_at, updated_at, filename FROM storage where key = ?`,
		key)
	return
}
func (s *dbStore) Set(data DBRecord) (err error) {
	_, err = s.db.Exec(`insert into storage 
 (key, description, created_at, updated_at, filename)
 values(?,?,DATETIME('now'),DATETIME('now'),?)
 on conflict (key) do update 
  set description=excluded.description,
      updated_at=excluded.updated_at,
      filename=excluded.filename`,
		data.Key, data.Description, data.Filename)
	return
}

func (s *dbStore) Delete(key string) (err error) {
	_, err = s.db.Exec(`delete from storage where key = ?`, key)
	return
}
