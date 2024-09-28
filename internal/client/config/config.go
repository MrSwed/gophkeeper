package config

import (
	"encoding/json"
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

type duration time.Duration

func (d duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

var (
	excludeSaveKeys  = []string{"config_path", "loaded_at", "changed_at", "encryption_key"}
	excludeViewKeys  = []string{"encryption_key"}
	durationViewKeys = []string{"timeout"}
	User             config
	Glob             = config{Viper: viper.New()}
)

func (c *config) Set(key string, value any) {
	c.Viper.Set(key, value)
	c.Viper.Set("changed_at", time.Now())
}

func (c *config) AllSettings() (m map[string]any) {
	m = c.Viper.AllSettings()
	for _, k := range excludeViewKeys {
		delete(m, k)
	}
	for _, k := range durationViewKeys {
		m[k] = duration(c.Viper.GetDuration(k))
	}
	return
}

func (c *config) Save() error {
	isNew := true
	if c.Viper.Get("loaded_at") != nil {
		isNew = false
	}
	clearAfterSave := []string{"changed_at"}
	if c.excluded == nil {
		c.excluded = make(map[string]any)
	}
	for _, k := range excludeSaveKeys {
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
		return c.Viper.SafeWriteConfig()
	}
	return c.Viper.WriteConfig()
}

func (c *config) IsChanged() bool {
	return c.Viper.Get("changed_at") != nil
}

func (c *config) Load(name, path string, defaults map[string]any) error {
	c.path = path
	c.Viper.SetConfigName(name)
	c.Viper.SetConfigType("json")
	c.Viper.AddConfigPath(path)
	c.Viper.AddConfigPath(".")
	for k, v := range defaults {
		c.Viper.SetDefault(k, v)
	}
	// _ = os.MkdirAll(path, 0755)

	err := c.Viper.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil
	}
	c.Viper.Set("loaded_at", time.Now())
	return err
}
