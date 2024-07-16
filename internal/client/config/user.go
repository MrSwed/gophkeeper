package config

import (
	"os"
	"path/filepath"
)

func UserLoad() (err error) {
	var name, cfgDir string
	if name = Glob.GetString("profile"); name == "" {
		name = "default"
	}
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		return
	}
	usrCfgDir := filepath.Join(cfgDir, AppName, name)
	profiles := Glob.GetStringMap("profiles")
	ch := false
	profile, ok := profiles[name].(map[string]any)
	if !ok {
		profile = newGlobProfileItem(usrCfgDir)
		ch = true
	}
	if _, ok = profile["path"].(string); !ok {
		profile["path"] = usrCfgDir
		ch = true
	}
	if err = User.Load(AppName, profile["path"].(string), map[string]any{"name": name, "autosave": true}); err == nil && ch {
		profiles[name] = profile
		Glob.Set("profiles", profiles)
	}
	return
}

func init() {
	err := UserLoad()
	if err != nil {
		panic(err)
	}
}
