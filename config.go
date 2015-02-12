package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	HostName                string
	ListenAt                string
	RethinkConnectionString string
	ForceHttps              bool
	JobInvertal             int
}

var (
	config *Config
)

func (c *Config) GetBaseUrl(host string) string {
	var h string
	if len(c.HostName) > 1 {
		h = c.HostName
	} else {
		h = host
	}
	if c.ForceHttps {
		return "https://" + h + "/"
	} else {
		return "http://" + h + "/"
	}
}

func MakeConfig() error {
	config = &Config{}
	if _, err := toml.DecodeFile("urls.toml", config); err != nil {
		return err
	} else {
		return nil
	}
}
