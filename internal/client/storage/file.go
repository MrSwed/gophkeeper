package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	if strings.Contains(fileName, "..") || filepath.IsAbs(fileName) {
		return nil, fmt.Errorf("invalid file name: %s", fileName)
	}
	fullPath := filepath.Join(s.path, fileName)
	fullPath = filepath.Clean(fullPath)
	relPath, err := filepath.Rel(s.path, fullPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return nil, fmt.Errorf("attempt to access file outside of allowed directory: %s", fullPath)
	}

	return os.ReadFile(fullPath)
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
	if strings.Contains(filePath, "..") || filepath.IsAbs(filePath) {
		return nil, fmt.Errorf("invalid file path: %s", filePath)
	}
	fullPath := filepath.Join(s.path, filePath)
	fullPath = filepath.Clean(fullPath)
	relPath, err := filepath.Rel(s.path, fullPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return nil, fmt.Errorf("attempt to access file outside of allowed directory: %s", fullPath)
	}

	return os.ReadFile(fullPath)
}

func (s *fileStore) SaveOrigin(filePath string, b []byte) (err error) {
	return os.WriteFile(filePath, b, 0600)
}
