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
}

var (
	User = config{viper.New()}
	Glob = config{viper.New()}
)

func (c *config) Print() {
	v := make(map[string]any)
	if err := c.Unmarshal(&v); err != nil {
		fmt.Println(err)
	}
	fmt.Printf(" %v\n", v)
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
	if isNew {
		return c.SafeWriteConfig()
	}
	return c.WriteConfig()
}

func (c *config) IsChanged() bool {
	return c.Get("changed_at") != nil
}

func (c *config) Load(name, path string, defaults map[string]any) error {
	c.SetConfigName(name)
	c.SetConfigType("json")
	c.AddConfigPath(path)
	c.AddConfigPath(".")

	for k, v := range defaults {
		c.Viper.Set(k, v)
	}
	_ = os.MkdirAll(path, 0755)

	err := c.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return nil
	}
	c.Viper.Set("loaded_at", time.Now())
	return err

}

func UserLoad(name string) (err error) {
	var cfgDir string
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		return
	}
	profiles, ok := Glob.Get("profiles").(map[string]any)
	if !ok {
		err = errors.New("error get profiles")
		return
	}
	profile, ok := profiles[name].(map[string]any)
	if !ok {
		err = errors.New("error get profile " + name)
		return
	}
	path, ok := profile["path"].(string)
	if !ok {
		profile["path"] = filepath.Join(cfgDir, name)
	}
	err = User.Load(AppName, path, map[string]any{"name": name, "autosave": true})
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
		"profile":  "default",
		"profiles": map[string]any{"default": map[string]any{
			"path": filepath.Join(cfgDir, "default"),
		}},
	})
	return
}

func init() {
	err := GlobalLoad()
	if err != nil {
		panic(err)
	}
	err = UserLoad(Glob.GetString("profile"))
	if err != nil {
		panic(err)
	}

}
