package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func GetUserName() (userName string) {
	userName = Glob.Viper.GetString("profile")
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
	if cfgDir = Glob.Viper.GetString("config_path"); cfgDir == "" {
		if cfgDir, err = os.UserConfigDir(); err != nil {
			return
		}
	}
	usrCfgDir = filepath.Join(cfgDir, AppName, userName)
	profiles := Glob.Viper.GetStringMap("profiles")
	profile, ok := profiles[strings.ToLower(userName)].(map[string]any)
	if !ok {
		profile = NewGlobProfileItem(usrCfgDir)
		ch = true
	}
	if _, ok = profile["name"].(string); !ok {
		profile["name"] = userName
	}
	if profilePath, ok := profile["path"].(string); !ok {
		profile["path"] = usrCfgDir
		ch = true
	} else {
		usrCfgDir = profilePath
	}
	if ch {
		profiles[strings.ToLower(userName)] = profile
		Glob.Set("profiles", profiles)
	}
	return
}

func UserLoad(reload ...bool) (err error) {
	err = GlobalLoad()
	if err != nil {
		return
	}
	userName := GetUserName()
	if User.Viper != nil && User.Viper.Get("loaded_at") != nil &&
		!(len(reload) > 0 && reload[0]) &&
		userName == User.Viper.GetString("name") {
		return
	}
	User = config{Viper: viper.New()}
	usrCfgDir := ""
	if usrCfgDir, err = UsrCfgDir(); err != nil {
		return
	}
	if err = os.MkdirAll(usrCfgDir, 0750); err != nil {
		return
	}
	err = User.Load(AppName, usrCfgDir,
		map[string]any{
			"name":                  userName,
			"autosave":              true,
			"db_file":               path.Join(usrCfgDir, userName+".db"),
			"sync.timeout.register": time.Minute * 1,
			"sync.timeout.sync":     time.Hour * 3,
		})
	return
}
