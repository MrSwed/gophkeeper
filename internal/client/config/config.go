package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/spf13/viper"
)

const (
	AppName = "GophKeeper"

	PageSize = 1000
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
	excludeSaveKeys          = []string{"config_path", "loaded_at", "changed_at", "encryption_key", "sync_password"}
	excludeViewKeys          = []string{"encryption_key", "sync_password"}
	durationViewKeys         = []string{"sync.timeout.sync", "sync.timeout.register"}
	clearAfterSave           = []string{"changed_at"}
	syncUpdatedTriggerFields = []string{"email", "packed_key", "sync.token"}
	User                     config
	Glob                     = config{Viper: viper.New()}
)

// Set
// own method with changed_at mark
func (c *config) Set(key string, value any) {
	c.Viper.Set(key, value)
	c.Viper.Set("changed_at", time.Now())
	for _, k := range syncUpdatedTriggerFields {
		if k == key {
			c.Viper.Set("sync.user.updated_at", time.Now())
			break
		}
	}
}

// deepMapSet
// deep correct value of config, usable at unmarshal,
// for example time.Duration as string
func deepMapSet(path []string, value any) (res map[string]any) {
	res = make(map[string]any)
	if len(path) >= 1 {
		if len(path) == 1 {
			res[path[0]] = value
			return
		}
		res[path[0]] = deepMapSet(path[1:], value)
		return
	}
	return
}

// AllSettings
// get all printable settings
func (c *config) AllSettings() (m map[string]any) {
	m = c.Viper.AllSettings()
	for _, k := range excludeViewKeys {
		delete(m, k)
	}
	for _, k := range durationViewKeys {
		path := strings.Split(k, ".")
		if len(path) == 1 {
			m[k] = duration(c.Viper.GetDuration(k))
		} else {
			upd := deepMapSet(path[:], duration(c.Viper.GetDuration(k)))
			_ = mergo.Map(&m, &upd, mergo.WithOverride)
		}
	}
	return
}

func (c *config) Save() error {
	isNew := true
	if c.Viper.Get("loaded_at") != nil {
		isNew = false
	}
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
	_ = os.MkdirAll(c.path, 0750)

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
