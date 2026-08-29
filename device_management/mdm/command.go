// Package mdm turns validated, generated payload structs into syntactically
// correct MDM artifacts: command plists and configuration profiles
// (.mobileconfig). It is the MDM half of the SDK's workflow engine — typed
// values in, spec-validated Apple config out.
package mdm

import (
	"crypto/sha256"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/plistenc"
)

// CommandPayload is implemented by every generated MDM command struct.
type CommandPayload interface {
	// RequestType is the command's wire identifier, e.g. "DeviceLock".
	RequestType() string
	// Validate checks the payload against Apple's spec.
	Validate() error
}

// CommandOption customizes command envelope construction.
type CommandOption func(*commandConfig)

type commandConfig struct {
	uuid string
}

// WithCommandUUID sets an explicit CommandUUID. Without it, a deterministic
// UUID is derived from the command content, so identical commands produce
// identical plists.
func WithCommandUUID(uuid string) CommandOption {
	return func(c *commandConfig) { c.uuid = uuid }
}

// NewCommand validates payload and renders the MDM command plist:
//
//	{ CommandUUID, Command: { RequestType, …payload keys } }
func NewCommand(payload CommandPayload, opts ...CommandOption) ([]byte, error) {
	if payload == nil {
		return nil, fmt.Errorf("mdm: nil command payload")
	}
	var cfg commandConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("mdm: invalid %s command: %w", payload.RequestType(), err)
	}

	fields, err := plistenc.Fields(payload)
	if err != nil {
		return nil, fmt.Errorf("mdm: %s: %w", payload.RequestType(), err)
	}
	body := append(plistenc.Dict{{Key: "RequestType", Value: payload.RequestType()}}, fields...)

	uuid := cfg.uuid
	if uuid == "" {
		doc, err := plistenc.Document(body)
		if err != nil {
			return nil, err
		}
		uuid = deriveUUID(payload.RequestType(), doc)
	}
	return plistenc.Document(plistenc.Dict{
		{Key: "CommandUUID", Value: uuid},
		{Key: "Command", Value: body},
	})
}

// deriveUUID builds a stable, UUID-shaped identifier from content, keeping
// command emission deterministic by default.
func deriveUUID(parts ...any) string {
	h := sha256.New()
	for _, p := range parts {
		fmt.Fprintf(h, "%v\x00", p)
	}
	s := h.Sum(nil)
	return fmt.Sprintf("%X-%X-%X-%X-%X", s[0:4], s[4:6], s[6:8], s[8:10], s[10:16])
}
