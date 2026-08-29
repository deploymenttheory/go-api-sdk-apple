// Build a configuration profile (.mobileconfig) with a restrictions
// payload. The profile envelope (PayloadContent, identifiers, UUIDs,
// versions) is assembled and validated by the SDK.
package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/profiles"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

func main() {
	doc, err := mdm.NewProfile("com.example.restrictions",
		mdm.WithDisplayName("Corporate Restrictions"),
		mdm.WithOrganization("Example Corp"),
		mdm.WithScope("System"),
		mdm.WithPayload(&profiles.Applicationaccess{
			AllowCamera:     ptr.To(false),
			AllowScreenShot: ptr.To(false),
			AllowAirDrop:    ptr.To(true),
		}),
	)
	if err != nil {
		log.Fatalf("build profile: %v", err)
	}
	fmt.Print(string(doc))
}
