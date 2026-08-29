// Command gendm is the codegen stage of the device-management SDK
// pipeline: it reads the committed spec snapshots (metadata/specs) and
// emits the generated payload structs, validation and registries under
// mdm/ and ddm/. It is offline and deterministic; CI regenerates and diffs
// against the committed tree.
//
//	go run ./device_management/cmd/gendm
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen"
)

func main() {
	metadataDir := flag.String("metadata", filepath.Join("device_management", "metadata", "specs"), "snapshot input directory")
	outDir := flag.String("out", "device_management", "output directory (the device_management root)")
	flag.Parse()

	if err := codegen.Run(*metadataDir, *outDir); err != nil {
		fmt.Fprintln(os.Stderr, "gendm:", err)
		os.Exit(1)
	}
}
