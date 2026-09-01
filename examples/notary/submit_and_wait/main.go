package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary"
)

func main() {
	fmt.Println("=== Apple Notary API - Submit and Wait ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "57246542-96fe-1a63-e053-0824d011072a"
	privateKeyPEM := `-----BEGIN PRIVATE KEY-----
your-app-store-connect-api-key
-----END PRIVATE KEY-----`

	privateKey, err := notary.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	c, err := notary.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Alternatively, read credentials from the environment:
	//   APPLE_KEY_ID, APPLE_ISSUER_ID, and APPLE_PRIVATE_KEY_PATH (or _PEM)
	// c, err := notary.NewClientFromEnv()

	ctx := context.Background()

	res, err := c.SubmitAndWait(ctx, notary.SubmitInput{
		FilePath: "./MyApp.pkg",
		Name:     "MyApp.pkg",
		// Webhook is optional: Apple posts the verdict to this URL when it
		// finishes, so a job need not sit and poll. It is best effort.
		// Webhook: "https://example.com/notary-callback",
		UploadProgress: func(written, total int64) {
			fmt.Printf("\ruploaded %d/%d bytes", written, total)
		},
	}, notary.WaitOptions{
		Interval: 30 * time.Second,
		Timeout:  30 * time.Minute,
		Progress: func(s *notary.Status) {
			fmt.Printf("\nstatus: %s\n", s.Status)
		},
	})

	// A submission that was not accepted returns ErrRejected with the log
	// issues already fetched and parsed.
	if errors.Is(err, notary.ErrRejected) {
		fmt.Printf("\nNotarization was not accepted: %s\n", res.Status.Status)
		for _, issue := range res.Issues {
			fmt.Printf("  [%s] %s (%s)\n", issue.Severity, issue.Message, issue.Path)
		}
		log.Fatalf("submission %s rejected", res.SubmissionID)
	}
	if err != nil {
		log.Fatalf("\nNotarization failed: %v", err)
	}

	fmt.Printf("\nAccepted. Submission ID: %s\n", res.SubmissionID)
	fmt.Println("Staple the ticket with: xcrun stapler staple ./MyApp.pkg")
}
