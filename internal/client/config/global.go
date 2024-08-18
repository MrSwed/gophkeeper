package config

import (
	"os"
	"path/filepath"
)

const (
	MaxBlobSize = 1024 * 64
)

func NewGlobProfileItem(path string) map[string]any {
	return map[string]any{
		"path": path,
	}
}

func GlobalLoad() (err error) {
	var cfgDir string
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		return
	}
	cfgDir = filepath.Join(cfgDir, AppName)
	err = Glob.Load(AppName, cfgDir, map[string]any{
		"autosave": true,
	})
	return
}
