package config

import (
	"time"

	"gophKeeper/internal/helper"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

// StorageConfig file storage configs
type StorageConfig struct {
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn" flag:"d" usage:"Provide the database dsn connect string"`
	// FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path" flag:"f" usage:"Provide the file storage path"`
}

type GRPC struct {
	GRPCAddress          string        `env:"GRPC_ADDRESS" json:"grpc_address"  flag:"g" usage:"Provide the grpc service address" envDefault:":3200"`
	GRPCOperationTimeout time.Duration `env:"GRPC_OPERATION_TIMEOUT" json:"grpc_operation_timeout" flag:"" usage:"Provide the grpc operation timeout" envDefault:"5s"`
}

// Config all configs
type Config struct {
	GRPC
	StorageConfig
}

func NewConfig() *Config {
	return &Config{}
}

// Init all configs
func (c *Config) Init() (*Config, error) {
	err := c.parseFlags()
	if err != nil {
		return nil, err
	}
	err = c.ParseEnv()
	if err != nil {
		return nil, err
	}

	return c, err
}

// ParseEnv gets ENV configs
func (c *Config) ParseEnv() error {
	return env.Parse(c)
}

func (c *Config) parseFlags() (err error) {
	err = helper.GenerateFlags(c, pflag.CommandLine)
	if err != nil {
		return
	}
	pflag.Parse()
	return
}
