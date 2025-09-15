package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		ServerAddress:  "localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}
}

func Load() *Config {
	cfg := NewConfig()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.ServerAddress = addr
	}
	if poll := os.Getenv("POLL_INTERVAL"); poll != "" {
		if seconds, err := strconv.Atoi(poll); err == nil {
			cfg.PollInterval = time.Duration(seconds) * time.Second
		}
	}
	if report := os.Getenv("REPORT_INTERVAL"); report != "" {
		if seconds, err := strconv.Atoi(report); err == nil {
			cfg.ReportInterval = time.Duration(seconds) * time.Second
		}
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address (e.g. http://localhost:8080)")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "poll interval in seconds")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "report interval in seconds")

	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)

	flag.Parse()

	return cfg
}
