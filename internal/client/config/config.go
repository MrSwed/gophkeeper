package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	gkErr "gophKeeper/internal/client/errors"

	"github.com/kirsle/configdir"

	"github.com/creasty/defaults"
)

const (
	AppName        = "GophKeeper"
	configFileName = "config.json"
)

type Global struct {
}

type Config struct {
	ServerAddress string
	ServerType    string        `json:"server_type" default:"grpc"`
	SyncInterval  time.Duration `json:"sync_interval" default:"10m"`
	LogFileName   string
	user          string
	configFile    string
}

func NewUserConfig(user string) (c *Config, err error) {
	c = &Config{
		user: user,
	}
	if err = defaults.Set(c); err != nil {
		return
	}

	if c.user == "" {
		err = gkErr.ErrUserNameInvalid
		return
	}

	configPath := configdir.LocalConfig(AppName, c.user)
	err = configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		return
	}
	c.configFile = filepath.Join(configPath, configFileName)
	if _, err = os.Stat(c.configFile); os.IsNotExist(err) {
		err = c.SaveConfig()
	} else if err == nil {
		err = c.LoadConfig()
	}
	return
}

func (c *Config) LoadConfig() (err error) {
	var fh *os.File
	if fh, err = os.Open(c.configFile); err != nil {
		return
	}
	defer func() { err = fh.Close() }()
	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&c)
	return
}

func (c *Config) SaveConfig() (err error) {
	var fh *os.File
	if fh, err = os.Create(c.configFile); err != nil {
		return
	}
	defer func() { err = fh.Close() }()
	encoder := json.NewEncoder(fh)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&c)
	return
}
