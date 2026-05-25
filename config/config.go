package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var config Config

func Load() *Config {
	var config Config
	_, err := toml.DecodeFile(filepath.Join("config", "crestbridge.toml"), &config)
	if err != nil {
		panic(err)
	}
	return &config
}
