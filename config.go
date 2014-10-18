package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	BaseURL    string
	ListenAt   string
	DBAddress  string
	DBPassword string
}

var (
	config *Config
)

func MakeConfig() error {
	config = &Config{}
	if _, err := toml.DecodeFile("urls.toml", config); err != nil {
		return err
	} else {
		return nil
	}
}
