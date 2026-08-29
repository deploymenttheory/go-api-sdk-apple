package device_management_test

// End-to-end smoke tests over the generated surface: typed payloads in,
// validated Apple config out — MDM command plists, configuration profiles
// and DDM declaration JSON.

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/activations"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/configurations"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/commands"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/profiles"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/validate"
)

func TestSmokeCommandPlist(t *testing.T) {
	doc, err := mdm.NewCommand(&commands.DeviceLock{
		Message: ptr.To("Locked by IT"),
		PIN:     ptr.To("123456"),
	})
	if err != nil {
		t.Fatal(err)
	}
	out := string(doc)
	for _, want := range []string{
		`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"`,
		"<key>CommandUUID</key>",
		"<key>RequestType</key>",
		"<string>DeviceLock</string>",
		"<key>Message</key>",
		"<string>Locked by IT</string>",
		"<key>PIN</key>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("command plist missing %q:\n%s", want, out)
		}
	}

	// Deterministic by default: same payload, same plist.
	doc2, _ := mdm.NewCommand(&commands.DeviceLock{Message: ptr.To("Locked by IT"), PIN: ptr.To("123456")})
	if string(doc) != string(doc2) {
		t.Error("command emission is not deterministic")
	}
}

func TestSmokeProfileMobileconfig(t *testing.T) {
	doc, err := mdm.NewProfile("com.example.dictionary",
		mdm.WithDisplayName("Managed Dictionary"),
		mdm.WithScope("System"),
		mdm.WithPayload(&profiles.Dictionary{}),
	)
	if err != nil {
		t.Fatal(err)
	}
	out := string(doc)
	for _, want := range []string{
		"<key>PayloadContent</key>",
		"<string>com.apple.Dictionary</string>",
		"<key>PayloadIdentifier</key>",
		"<string>com.example.dictionary.0</string>",
		"<string>Configuration</string>",
		"<key>PayloadUUID</key>",
		"<key>PayloadDisplayName</key>",
		"<string>Managed Dictionary</string>",
		"<key>PayloadScope</key>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("profile missing %q:\n%s", want, out)
		}
	}
}

func TestSmokeDeclarationJSON(t *testing.T) {
	doc, err := ddm.BuildDeclaration("com.example.passcode", &configurations.PasscodeSettings{
		RequirePasscode: ptr.To(true),
		MinimumLength:   ptr.To(int64(8)),
	}, ddm.WithServerToken("v1"))
	if err != nil {
		t.Fatal(err)
	}
	var decl struct {
		Type        string
		Identifier  string
		ServerToken string
		Payload     map[string]any
	}
	if err := json.Unmarshal(doc, &decl); err != nil {
		t.Fatalf("declaration is not valid JSON: %v\n%s", err, doc)
	}
	if decl.Type != "com.apple.configuration.passcode.settings" || decl.Identifier != "com.example.passcode" || decl.ServerToken != "v1" {
		t.Fatalf("envelope = %+v", decl)
	}
	if decl.Payload["MinimumLength"] != float64(8) || decl.Payload["RequirePasscode"] != true {
		t.Fatalf("payload = %+v", decl.Payload)
	}
	// Optional keys the caller didn't set must be absent, not zeroed.
	if _, present := decl.Payload["RequireComplexPasscode"]; present {
		t.Error("unset optional key serialized")
	}
}

func TestSmokeValidationRejectsBadConfig(t *testing.T) {
	// Spec: MinimumLength has range 0..16.
	_, err := ddm.BuildDeclaration("com.example.passcode", &configurations.PasscodeSettings{
		MinimumLength: ptr.To(int64(99)),
	})
	if err == nil || !strings.Contains(err.Error(), "MinimumLength") {
		t.Fatalf("expected range violation, got %v", err)
	}

	// Spec: Activation:Simple requires StandardConfigurations.
	_, err = ddm.BuildDeclaration("com.example.activation", &activations.Simple{})
	if err == nil || !strings.Contains(err.Error(), "StandardConfigurations") {
		t.Fatalf("expected required-key violation, got %v", err)
	}
}

func TestSmokeRegistries(t *testing.T) {
	f, ok := commands.ByRequestType["DeviceLock"]
	if !ok {
		t.Fatal("DeviceLock missing from command registry")
	}
	if f().RequestType() != "DeviceLock" {
		t.Fatal("command factory returns wrong type")
	}
	if _, ok := profiles.ByPayloadType["com.apple.Dictionary"]; !ok {
		t.Fatal("com.apple.Dictionary missing from profile registry")
	}
	if _, ok := configurations.ByDeclarationType["com.apple.configuration.passcode.settings"]; !ok {
		t.Fatal("passcode settings missing from declaration registry")
	}
}

func TestSmokeValidateHelpersComposable(t *testing.T) {
	err := errors.Join(
		validate.InList("Mode", "bogus", []string{"auto", "manual"}),
		validate.IntRange("N", 99, ptr.To(int64(0)), ptr.To(int64(10))),
	)
	if err == nil || !strings.Contains(err.Error(), "Mode") || !strings.Contains(err.Error(), "N") {
		t.Fatalf("joined validation errors = %v", err)
	}
}
