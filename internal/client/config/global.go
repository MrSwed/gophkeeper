package config

import (
	"os"
	"path/filepath"
)

const (
	MaxBlobSize = 1024 * 64
)

type globProfileItem map[string]any

func newGlobProfileItem(path string) globProfileItem {
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
func init() {
	err := GlobalLoad()
	if err != nil {
		panic(err)
	}

}
