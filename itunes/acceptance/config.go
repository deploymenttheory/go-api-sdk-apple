package acceptance

import (
	"os"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/itunes"
)

// TestConfig holds configuration for iTunes acceptance tests, driven by
// environment variables. No credentials are required — the iTunes Search API
// is a public, unauthenticated API.
type TestConfig struct {
	RequestTimeout time.Duration
	Verbose        bool
}

var (
	// Config is the global acceptance test configuration, initialised from env.
	Config *TestConfig
	// Client is the shared iTunes SDK client for acceptance tests.
	Client *itunes.Client
)

func init() {
	Config = &TestConfig{
		RequestTimeout: getDurationEnv("ITUNES_REQUEST_TIMEOUT", 30*time.Second),
		Verbose:        getBoolEnv("ITUNES_VERBOSE", false),
	}
}

// InitClient creates and stores the shared iTunes client.
func InitClient() error {
	var err error
	Client, err = itunes.NewClient(
		itunes.WithTimeout(Config.RequestTimeout),
	)
	return err
}

// getDurationEnv parses a duration from an environment variable, returning
// fallback when the variable is absent or unparseable.
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

// getBoolEnv parses a boolean from an environment variable.
func getBoolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
