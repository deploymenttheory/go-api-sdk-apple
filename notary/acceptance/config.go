package acceptance

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary"
)

// TestConfig holds configuration for acceptance tests, driven by environment variables.
type TestConfig struct {
	// Auth — mirrors the variables read by notary.NewClientFromEnv.
	// Supply exactly one of PrivateKeyPEM or PrivateKeyPath.
	KeyID          string
	IssuerID       string
	PrivateKeyPEM  string // APPLE_PRIVATE_KEY_PEM  — inline PEM content
	PrivateKeyPath string // APPLE_PRIVATE_KEY_PATH — path to .p8 file

	// Test behaviour
	RequestTimeout time.Duration
	SkipCleanup    bool
	Verbose        bool
}

var (
	// Config is the global acceptance test configuration, initialised from env.
	Config *TestConfig
	// Client is the shared Notary SDK client for acceptance tests.
	Client *notary.Client
)

func init() {
	Config = &TestConfig{
		KeyID:          os.Getenv("APPLE_KEY_ID"),
		IssuerID:       os.Getenv("APPLE_ISSUER_ID"),
		PrivateKeyPEM:  os.Getenv("APPLE_PRIVATE_KEY_PEM"),
		PrivateKeyPath: os.Getenv("APPLE_PRIVATE_KEY_PATH"),
		RequestTimeout: getDurationEnv("NOTARY_REQUEST_TIMEOUT", 30*time.Second),
		SkipCleanup:    getBoolEnv("NOTARY_SKIP_CLEANUP", false),
		Verbose:        getBoolEnv("NOTARY_VERBOSE", false),
	}
}

// IsConfigured returns true when the minimum required credentials are present.
func IsConfigured() bool {
	return Config.KeyID != "" && Config.IssuerID != "" &&
		(Config.PrivateKeyPEM != "" || Config.PrivateKeyPath != "")
}

// InitClient creates and stores the shared Notary client from environment variables.
func InitClient() error {
	var privateKey any
	var err error

	switch {
	case Config.PrivateKeyPEM != "":
		privateKey, err = notary.ParsePrivateKey([]byte(Config.PrivateKeyPEM))
		if err != nil {
			return fmt.Errorf("failed to parse APPLE_PRIVATE_KEY_PEM: %w", err)
		}
	case Config.PrivateKeyPath != "":
		privateKey, err = notary.LoadPrivateKeyFromFile(Config.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key from %q: %w", Config.PrivateKeyPath, err)
		}
	default:
		return fmt.Errorf("either APPLE_PRIVATE_KEY_PEM or APPLE_PRIVATE_KEY_PATH must be set")
	}

	Client, err = notary.NewClient(
		Config.KeyID,
		Config.IssuerID,
		privateKey,
		notary.WithTimeout(Config.RequestTimeout),
	)
	if err != nil {
		return fmt.Errorf("failed to create Notary client: %w", err)
	}

	if Config.Verbose {
		log.Printf("Notary acceptance test client initialised (issuer: %s)", Config.IssuerID)
	}
	return nil
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
