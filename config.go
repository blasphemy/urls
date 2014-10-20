package main

import (
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	BaseURL     string
	ListenAt    string
	DBAddress   string
	DBPassword  string
	RunJobs     bool
	JobInvertal int
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

func (c *Config) GetJobInvertal() time.Duration {
	if c.JobInvertal > 0 {
		return time.Minute * time.Duration(c.JobInvertal)
	} else {
		return time.Minute * 10
	}
}
