package config

import "time"

type Config struct {
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		ServerAddress:  ":8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}
}
