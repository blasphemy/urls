package main

import (
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HostName                string
	ListenAt                string
	DBAddress               string
	DBPassword              string
	RethinkConnectionString string
	RunJobs                 bool
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

func (c *Config) GetJobInvertal() time.Duration {
	if c.JobInvertal > 0 {
		return time.Minute * time.Duration(c.JobInvertal)
	} else {
		return time.Minute * 10
	}
}
