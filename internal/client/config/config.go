package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	AppName = "GophKeeper"
)

type config struct {
	*viper.Viper
	path     string
	excluded []string
}

var (
	User = config{Viper: viper.New()}
	Glob = config{Viper: viper.New()}
)

func (c *config) Print() {
	// v := make(map[string]any)
	// if err := c.Unmarshal(&v); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%s %v\n", c.path, v)
	fmt.Println(c.AllSettings())
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
	excluded := append([]string{"loaded_at", "changed_at", "secret_pass"}, c.excluded...)
	for _, k := range excluded {
		c.Viper.Set(k, nil)
	}
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
