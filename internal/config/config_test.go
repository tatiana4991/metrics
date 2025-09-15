package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
}

func TestLoad_DefaultValues(t *testing.T) {
	resetFlags()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	cfg := Load()

	assert.Equal(t, "localhost:8080", cfg.ServerAddress)
	assert.Equal(t, 2*time.Second, cfg.PollInterval)
	assert.Equal(t, 10*time.Second, cfg.ReportInterval)
}

func TestLoad_EnvVarsOverrideDefaults(t *testing.T) {
	resetFlags()
	os.Setenv("ADDRESS", "http://test:9090")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("REPORT_INTERVAL", "15")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("REPORT_INTERVAL")
	}()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	cfg := Load()

	assert.Equal(t, "http://test:9090", cfg.ServerAddress)
	assert.Equal(t, 5*time.Second, cfg.PollInterval)
	assert.Equal(t, 15*time.Second, cfg.ReportInterval)
}

func TestLoad_FlagsOverrideEnvVars(t *testing.T) {
	resetFlags()
	os.Setenv("ADDRESS", "http://env:8080")
	os.Setenv("POLL_INTERVAL", "10")
	os.Setenv("REPORT_INTERVAL", "20")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("REPORT_INTERVAL")
	}()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	os.Args = []string{
		"cmd",
		"-a", "http://flag:9090",
		"-p", "3s",
		"-r", "25s",
	}

	cfg := Load()

	assert.Equal(t, "http://flag:9090", cfg.ServerAddress)
	assert.Equal(t, 3*time.Second, cfg.PollInterval)
	assert.Equal(t, 25*time.Second, cfg.ReportInterval)
}

func TestLoad_InvalidEnvValues(t *testing.T) {
	resetFlags()
	os.Setenv("POLL_INTERVAL", "invalid")
	os.Setenv("REPORT_INTERVAL", "not-a-number")
	defer func() {
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("REPORT_INTERVAL")
	}()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	cfg := Load()

	assert.Equal(t, 2*time.Second, cfg.PollInterval)
	assert.Equal(t, 10*time.Second, cfg.ReportInterval)
}

func TestLoad_FlagOverridesInvalidEnv(t *testing.T) {
	resetFlags()
	os.Setenv("POLL_INTERVAL", "invalid")
	defer os.Unsetenv("POLL_INTERVAL")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-p", "7s"}

	cfg := Load()

	assert.Equal(t, 7*time.Second, cfg.PollInterval)
}
