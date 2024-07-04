package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	AppName = "GophKeeper"
)

type config struct {
	*viper.Viper
	path string
}

var (
	User = config{Viper: viper.New()}
	Glob = config{Viper: viper.New()}
)

func (c *config) Print() {
	v := make(map[string]any)
	if err := c.Unmarshal(&v); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s %v\n", c.path, v)
}

func (c *config) Set(key string, value any) {
	c.Viper.Set(key, value)
	c.Viper.Set("changed_at", time.Now())
}

func (c *config) Save() error {
	isNew := true
	if c.Get("loaded_at") != nil {
		isNew = false
	}
	c.Viper.Set("loaded_at", nil)
	c.Viper.Set("changed_at", nil)
	c.Viper.Set("updated_at", time.Now())
	_ = os.MkdirAll(c.path, 0755)

	if isNew {
		return c.SafeWriteConfig()
	}
	return c.WriteConfig()
}

func (c *config) IsChanged() bool {
	return c.Get("changed_at") != nil
}

func (c *config) Load(name, path string, defaults map[string]any) error {
	c.path = path
	c.SetConfigName(name)
	c.SetConfigType("json")
	c.AddConfigPath(path)
	c.AddConfigPath(".")
	for k, v := range defaults {
		c.Viper.SetDefault(k, v)
	}
	// _ = os.MkdirAll(path, 0755)

	err := c.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil
	}
	c.Viper.Set("loaded_at", time.Now())
	return err

}

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
	err = UserLoad()
	if err != nil {
		panic(err)
	}

}
