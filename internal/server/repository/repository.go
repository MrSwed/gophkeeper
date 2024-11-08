package repository

import (
	"context"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	sqrl "github.com/Masterminds/squirrel"
)

var sq = sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)

const (
	storeTableName  = "storage"
	userTableName   = "users"
	clientTableName = "clients"
)

type store struct {
	db *sqlx.DB
	c  *config.StorageConfig
}

// DataStorage methods
type DataStorage interface {
	GetDataItem(ctx context.Context, userID uuid.UUID, key string) (item model.DBRecord, err error)
	ListDataItems(ctx context.Context, q *model.ListQuery) (item []model.ItemShort, err error)
	CountDataItems(ctx context.Context, q *model.ListQuery) (count uint64, err error)
	SaveDataItem(ctx context.Context, item model.DBRecord) (err error)
}

type UserStorage interface {
	GetUserSelf(ctx context.Context) (user model.DBUser, err error)
	GetUserByEmail(ctx context.Context, email string) (user model.DBUser, err error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (err error)
	GetUserIDByToken(ctx context.Context, token []byte) (userID uuid.UUID, err error)
	SaveUser(ctx context.Context, user *model.DBUser) (err error)
	NewUserClientToken(ctx context.Context, userID uuid.UUID, expAt *time.Time, meta any) (token []byte, err error)
}

// type FileStorage interface {
// 	GetFile(ctx context.Context, path string) ([]byte, error)
// 	SaveFile(ctx context.Context, path string, data []byte) error
// 	DeleteFile(ctx context.Context, path string) error
// }

type Storage interface {
	DataStorage
	UserStorage
	// FileStorage
}

type storage struct {
	DataStorage
	UserStorage
	// FileStorage
}

var _ Storage = (*storage)(nil)

// NewRepository return repository of database or memory if no db set
func NewRepository(c *config.StorageConfig, db *sqlx.DB) (s Storage) {
	return &storage{
		DataStorage: NewDBRepository(c, db),
		UserStorage: NewUserStorage(c, db),
		// FileStorage: NewFileStorageRepository(c),
	}
}
