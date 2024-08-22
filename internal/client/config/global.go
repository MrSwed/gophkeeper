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

func GlobalLoad(reload ...bool) (err error) {
	if Glob.Get("loaded_at") != nil && (len(reload) == 0 || !reload[0]) {
		return
	}
	var cfgDir string
	// can set configs home path before load
	if Glob.GetString("config_path") == "" {
		cfgDir, err = os.UserConfigDir()
		if err != nil {
			return
		}
		Glob.Viper.Set("config_path", cfgDir)
	} else {
		cfgDir = Glob.GetString("config_path")
	}
	cfgDir = filepath.Join(cfgDir, AppName)
	err = Glob.Load(AppName, cfgDir, map[string]any{
		"autosave": true,
	})
	return
}
