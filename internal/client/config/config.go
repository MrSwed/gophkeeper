package config

import (
	"errors"
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
	excluded map[string]any
}

var (
	User config
	Glob = config{Viper: viper.New()}
)

func (c *config) Set(key string, value any) {
	c.Viper.Set(key, value)
	c.Viper.Set("changed_at", time.Now())
}

func (c *config) Save() error {
	isNew := true
	if c.Get("loaded_at") != nil {
		isNew = false
	}
	clearAfterSave := []string{"changed_at"}
	excluded := []string{"config_path", "loaded_at", "changed_at", "encryption_key"}
	if c.excluded == nil {
		c.excluded = make(map[string]any)
	}
	for _, k := range excluded {
		if x := c.Viper.Get(k); x != nil {
			c.excluded[k] = x
			c.Viper.Set(k, nil)
			for _, clr := range clearAfterSave {
				if clr == k {
					c.excluded[k] = nil
				}
			}
		}
	}
	defer func() {
		// restore excluded fields
		// todo: check viper for excluded from save instead
		for k, v := range c.excluded {
			if v != nil {
				c.Viper.Set(k, v)
			}
		}
	}()

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
