package config

import "time"

type Config struct {
	ServerAddress string
	ServerType    string        `json:"server_type" default:"http"`
	SyncInterval  time.Duration `json:"sync_interval" default:"10m"`
	LogFileName   string
}
