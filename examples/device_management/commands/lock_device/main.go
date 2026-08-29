// Build a DeviceLock MDM command plist from typed, spec-validated values.
package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/commands"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

func main() {
	doc, err := mdm.NewCommand(&commands.DeviceLock{
		Message:     ptr.To("This Mac has been locked by IT."),
		PhoneNumber: ptr.To("+44 20 7946 0000"),
		PIN:         ptr.To("123456"),
	})
	if err != nil {
		log.Fatalf("build command: %v", err)
	}
	fmt.Print(string(doc))
}
