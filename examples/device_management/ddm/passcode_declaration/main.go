// Build a DDM passcode-policy declaration as JSON. Validation enforces the
// spec: try MinimumLength: 99 and the build fails (allowed range is 0-16).
package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/configurations"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

func main() {
	doc, err := ddm.BuildDeclaration("com.example.passcode-policy",
		&configurations.PasscodeSettings{
			RequirePasscode:          ptr.To(true),
			RequireComplexPasscode:   ptr.To(true),
			MinimumLength:            ptr.To(int64(12)),
			MaximumFailedAttempts:    ptr.To(int64(6)),
			MaximumPasscodeAgeInDays: ptr.To(int64(365)),
		},
		ddm.WithServerToken("2026-07-17.1"),
	)
	if err != nil {
		log.Fatalf("build declaration: %v", err)
	}
	fmt.Print(string(doc))
}
