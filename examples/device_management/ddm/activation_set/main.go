// Build a DDM activation that ties a set of configurations together — the
// complete declaration set a DDM server would serve for a passcode policy.
package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/activations"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/configurations"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

func main() {
	config, err := ddm.BuildDeclaration("com.example.passcode-policy",
		&configurations.PasscodeSettings{MinimumLength: ptr.To(int64(8))})
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	activation, err := ddm.BuildDeclaration("com.example.passcode-activation",
		&activations.Simple{
			StandardConfigurations: []string{"com.example.passcode-policy"},
		})
	if err != nil {
		log.Fatalf("activation: %v", err)
	}

	fmt.Println("— configuration —")
	fmt.Print(string(config))
	fmt.Println("— activation —")
	fmt.Print(string(activation))
}
