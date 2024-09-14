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

type store struct {
	db *sqlx.DB
	c  *config.StorageConfig
}

// DataStorage methods
type DataStorage interface {
	GetItem(ctx context.Context, id string) (item *model.DBRecord, err error)
	ListItems(ctx context.Context, q *model.ListQuery) (item []model.ItemShort, err error)
	CountItems(ctx context.Context, q *model.ListQuery) (count int, err error)
	SaveItem(ctx context.Context, item *model.DBRecord) (err error)
	DeleteItem(ctx context.Context, id string) (err error)
}

type UserStorage interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (user model.User, err error)
	GetByEmail(ctx context.Context, email string) (user model.User, err error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (err error)
	DeleteClient(ctx context.Context, token []byte) (err error)
	GetUserIDByToken(ctx context.Context, token []byte) (userID uuid.UUID, err error)
	NewUser(ctx context.Context, user model.User) (userID uuid.UUID, err error)
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
func NewRepository(c *config.StorageConfig, db *sqlx.DB) (s *storage) {
	return &storage{
		DataStorage: NewDBRepository(c, db),
		UserStorage: NewAuthStorage(c, db),
		// FileStorage: NewFileStorageRepository(c),
	}
}
