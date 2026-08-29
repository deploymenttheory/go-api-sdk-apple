package main

import (
	"context"
	"fmt"
	"log"

	microsoft_updates "github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates"
)

func main() {
	c, err := microsoft_updates.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.MicrosoftUpdatesAPI.OneDrive.GetAllRingsV1(context.Background())
	if err != nil {
		log.Fatalf("OneDrive.GetAllRingsV1: %v", err)
	}

	fmt.Printf("OneDrive distribution rings (%d found)\n", len(resp.Rings))
	fmt.Println("=========================================")
	for _, ring := range resp.Rings {
		fmt.Printf("%-20s version=%-20s url=%s\n",
			ring.Ring,
			ring.Version,
			ring.DownloadURL,
		)
	}
}
