package repository

/* * /

import (
	"context"
	"gophKeeper/internal/helper"
	"gophKeeper/internal/server/config"
	"os"

	"github.com/google/uuid"
)

func NewFileStorageRepository(c *config.StorageConfig) *fileStore {
	return &fileStore{
		c: c,
	}
}

// todo maybe feature use
// var _ FileStorage = (*fileStore)(nil)

type fileStore struct {
	c *config.StorageConfig
}

func (f *fileStore) GetFile(ctx context.Context, path string) (blob []byte, err error) {
	var userID uuid.UUID
	userID, err = helper.GetCtxUserID(ctx)
	// TODO implement me
	panic("implement me " + userID.String())
}

func (f *fileStore) SaveFile(ctx context.Context, path string, data []byte) (err error) {
	var userID uuid.UUID
	userID, err = helper.GetCtxUserID(ctx)
	// TODO implement me
	panic("implement me " + userID.String())
}
func (f *fileStore) DeleteFile(_ context.Context, fileName string) (err error) {
	return os.Remove(fileName)
}

/**/
