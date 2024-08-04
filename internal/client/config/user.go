package config

import (
	"os"
	"path"
	"path/filepath"
)

func GetUserName() (userName string) {
	userName = Glob.GetString("profile")
	if userName == "" {
		userName = "default"
	}
	return
}

func UsrCfgDir(userNames ...string) (p string, err error) {
	var userName, cfgDir string
	if len(userNames) > 0 {
		userName = userNames[0]
	} else {
		userName = GetUserName()
	}
	cfgDir, err = os.UserConfigDir()
	p = filepath.Join(cfgDir, AppName, userName)
	return
}

func UserLoad() (err error) {
	var (
		userName, usrCfgDir string
	)
	userName = GetUserName()

	if usrCfgDir, err = UsrCfgDir(); err != nil {
		return
	}

	if err = os.MkdirAll(usrCfgDir, os.ModePerm); err != nil {
		return
	}
	profiles := Glob.GetStringMap("profiles")
	ch := false
	profile, ok := profiles[userName].(map[string]any)
	if !ok {
		profile = newGlobProfileItem(usrCfgDir)
		ch = true
	}
	if _, ok = profile["path"].(string); !ok {
		profile["path"] = usrCfgDir
		ch = true
	}

	if err = User.Load(AppName, profile["path"].(string),
		map[string]any{
			"name":     userName,
			"autosave": true,
			"db_file":  path.Join(usrCfgDir, userName+".db"),
		}); err == nil && ch {
		profiles[userName] = profile
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
