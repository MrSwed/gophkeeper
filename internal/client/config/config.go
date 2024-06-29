package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kirsle/configdir"
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

func UserLoad(name string) error {
	return User.Load(AppName+"_"+name, configdir.LocalConfig(AppName, name),
		map[string]any{"name": name, "autosave": true})
}

func GlobalLoad() error {
	return Glob.Load(AppName, configdir.LocalConfig(AppName), map[string]any{"autosave": true})
}

func init() {
	// UserLoad()
	err := GlobalLoad()
	if err != nil {
		panic(err)
	}
}
