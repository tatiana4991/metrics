package main

import (
	"flag"
	"time"

	"github.com/tatiana4991/metrics/internal/agent/collector"
	"github.com/tatiana4991/metrics/internal/agent/sender"
	"github.com/tatiana4991/metrics/internal/config"
	"github.com/tatiana4991/metrics/internal/storage"
)

func main() {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.ServerAddress, "a", "http://localhost:8080", "Server address (default http://localhost:8080)")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "Poll interval in seconds (default 2s)")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval in seconds (default 10s)")
	flag.Parse()

	store := storage.NewMemStorage()
	collector := collector.NewCollector(store)
	sender := sender.NewSender(cfg, store)

	go func() {
		ticker := time.NewTicker(cfg.PollInterval)
		defer ticker.Stop()
		for range ticker.C {
			collector.Collect()
		}
	}()

	ticker := time.NewTicker(cfg.ReportInterval)
	defer ticker.Stop()
	for range ticker.C {
		sender.SendAll()
	}
}
