package main

import (
	"time"

	"github.com/tatiana4991/metrics/internal/agent/collector"
	"github.com/tatiana4991/metrics/internal/agent/sender"
	"github.com/tatiana4991/metrics/internal/config"
	"github.com/tatiana4991/metrics/internal/storage"
)

func main() {
	cfg := config.Load()
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
