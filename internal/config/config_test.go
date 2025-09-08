package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, ":8080", cfg.ServerAddress)
	assert.Equal(t, 2*time.Second, cfg.PollInterval)
	assert.Equal(t, 10*time.Second, cfg.ReportInterval)
}
