package storage

import (
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type fileStore struct {
	path string
}

func NewFileStore(path string) *fileStore {
	return &fileStore{
		path: path,
	}
}

func (s *fileStore) GetStored(fileName string) (b []byte, err error) {
	return os.ReadFile(filepath.Join(s.path, fileName))
}

func (s *fileStore) SaveStore(fileName string, b []byte) (err error) {
	if err = os.Mkdir(s.path, 0600); err != nil && !os.IsExist(err) {
		return
	}
	return os.WriteFile(filepath.Join(s.path, fileName), b, 0600)
}

func (s *fileStore) Delete(fileName string) (err error) {
	return os.Remove(filepath.Join(s.path, fileName))
}

func (s *fileStore) GetOrigin(filePath string) (b []byte, err error) {
	return os.ReadFile(filePath)
}
func (s *fileStore) SaveOrigin(filePath string, b []byte) (err error) {
	return os.WriteFile(filePath, b, 0644)
}
