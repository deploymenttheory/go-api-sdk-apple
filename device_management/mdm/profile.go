package mdm

import (
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/plistenc"
)

// ProfilePayload is implemented by every generated profile payload struct.
type ProfilePayload interface {
	// PayloadType is the payload's wire identifier, e.g. "com.apple.wifi.managed".
	PayloadType() string
	// Validate checks the payload against Apple's spec.
	Validate() error
}

// ProfileOption customizes profile construction.
type ProfileOption func(*profileConfig)

type profileConfig struct {
	displayName          string
	description          string
	organization         string
	scope                string
	uuid                 string
	removalDisallowed    *bool
	durationUntilRemoval *int64
	payloads             []ProfilePayload
}

// WithDisplayName sets PayloadDisplayName on the profile.
func WithDisplayName(name string) ProfileOption {
	return func(c *profileConfig) { c.displayName = name }
}

// WithDescription sets PayloadDescription on the profile.
func WithDescription(desc string) ProfileOption {
	return func(c *profileConfig) { c.description = desc }
}

// WithOrganization sets PayloadOrganization on the profile.
func WithOrganization(org string) ProfileOption {
	return func(c *profileConfig) { c.organization = org }
}

// WithScope sets PayloadScope: "System" or "User".
func WithScope(scope string) ProfileOption {
	return func(c *profileConfig) { c.scope = scope }
}

// WithProfileUUID sets an explicit PayloadUUID for the profile. Without it,
// a deterministic UUID is derived from the profile identifier.
func WithProfileUUID(uuid string) ProfileOption {
	return func(c *profileConfig) { c.uuid = uuid }
}

// WithRemovalDisallowed sets PayloadRemovalDisallowed.
func WithRemovalDisallowed(disallowed bool) ProfileOption {
	return func(c *profileConfig) { c.removalDisallowed = &disallowed }
}

// WithPayload appends a payload to the profile's PayloadContent.
func WithPayload(p ProfilePayload) ProfileOption {
	return func(c *profileConfig) { c.payloads = append(c.payloads, p) }
}

// NewProfile validates every payload and renders a configuration profile
// (.mobileconfig plist). identifier is the profile's PayloadIdentifier;
// each payload's PayloadIdentifier and PayloadUUID are derived from it
// deterministically (identifier.<n>).
func NewProfile(identifier string, opts ...ProfileOption) ([]byte, error) {
	if identifier == "" {
		return nil, fmt.Errorf("mdm: profile identifier is required")
	}
	var cfg profileConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	if len(cfg.payloads) == 0 {
		return nil, fmt.Errorf("mdm: profile %s has no payloads", identifier)
	}

	var content []any
	for i, p := range cfg.payloads {
		if p == nil {
			return nil, fmt.Errorf("mdm: profile %s: payload %d is nil", identifier, i)
		}
		if err := p.Validate(); err != nil {
			return nil, fmt.Errorf("mdm: profile %s: invalid %s payload: %w", identifier, p.PayloadType(), err)
		}
		fields, err := plistenc.Fields(p)
		if err != nil {
			return nil, fmt.Errorf("mdm: profile %s: %s: %w", identifier, p.PayloadType(), err)
		}
		payloadID := fmt.Sprintf("%s.%d", identifier, i)
		entry := plistenc.Dict{
			{Key: "PayloadType", Value: p.PayloadType()},
			{Key: "PayloadIdentifier", Value: payloadID},
			{Key: "PayloadUUID", Value: deriveUUID(payloadID, p.PayloadType())},
			{Key: "PayloadVersion", Value: int64(1)},
		}
		content = append(content, append(entry, fields...))
	}

	uuid := cfg.uuid
	if uuid == "" {
		uuid = deriveUUID(identifier)
	}
	root := plistenc.Dict{
		{Key: "PayloadContent", Value: content},
		{Key: "PayloadIdentifier", Value: identifier},
		{Key: "PayloadType", Value: "Configuration"},
		{Key: "PayloadUUID", Value: uuid},
		{Key: "PayloadVersion", Value: int64(1)},
	}
	if cfg.displayName != "" {
		root = append(root, plistenc.Entry{Key: "PayloadDisplayName", Value: cfg.displayName})
	}
	if cfg.description != "" {
		root = append(root, plistenc.Entry{Key: "PayloadDescription", Value: cfg.description})
	}
	if cfg.organization != "" {
		root = append(root, plistenc.Entry{Key: "PayloadOrganization", Value: cfg.organization})
	}
	if cfg.scope != "" {
		root = append(root, plistenc.Entry{Key: "PayloadScope", Value: cfg.scope})
	}
	if cfg.removalDisallowed != nil {
		root = append(root, plistenc.Entry{Key: "PayloadRemovalDisallowed", Value: *cfg.removalDisallowed})
	}
	return plistenc.Document(root)
}
