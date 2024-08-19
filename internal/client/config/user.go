package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetUserName() (userName string) {
	userName = Glob.GetString("profile")
	if userName == "" {
		userName = "default"
		Glob.Set("profile", userName)
	}
	return
}

func UsrCfgDir(userNames ...string) (usrCfgDir string, err error) {
	var (
		userName, cfgDir string
		ch               bool
	)
	if len(userNames) > 0 {
		userName = userNames[0]
	} else {
		userName = GetUserName()
	}
	if cfgDir = Glob.GetString("config_path"); cfgDir == "" {
		if cfgDir, err = os.UserConfigDir(); err != nil {
			return
		}
	}
	usrCfgDir = filepath.Join(cfgDir, AppName, userName)
	profiles := Glob.GetStringMap("profiles")
	profile, ok := profiles[userName].(map[string]any)
	if !ok {
		profile = NewGlobProfileItem(usrCfgDir)
		ch = true
	}
	if profilePath, ok := profile["path"].(string); !ok {
		profile["path"] = usrCfgDir
		ch = true
	} else {
		usrCfgDir = profilePath
	}
	if ch {
		profiles[userName] = profile
		Glob.Set("profiles", profiles)
	}
	return
}

func UserLoad() (err error) {
	User = config{Viper: viper.New()}
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

	err = User.Load(AppName, usrCfgDir,
		map[string]any{
			"name":     userName,
			"autosave": true,
			"db_file":  path.Join(usrCfgDir, userName+".db"),
		})
	return
}
