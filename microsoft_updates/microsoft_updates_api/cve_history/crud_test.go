package cve_history_test

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/cve_history"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/cve_history/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) *cve_history.CVEHistoryService {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return cve_history.NewService(transport)
}

func TestGetCVEHistoryV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterCVEHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetCVEHistoryV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Entries)
}

func TestGetCVEHistoryV1_CVEExtraction(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterCVEHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetCVEHistoryV1(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Entries)

	// Find any entry that has CVEs.
	var entryWithCVEs *cve_history.CVEEntry
	for i := range resp.Entries {
		if len(resp.Entries[i].CVEs) > 0 {
			entryWithCVEs = &resp.Entries[i]
			break
		}
	}
	require.NotNil(t, entryWithCVEs, "expected at least one entry with CVEs")
	assert.NotEmpty(t, entryWithCVEs.CVEs)

	// All CVEs should match the CVE pattern.
	for _, cve := range entryWithCVEs.CVEs {
		assert.Regexp(t, `^CVE-\d{4}-\d{4,}$`, cve)
	}
}

func TestGetCVEHistoryV1_MockCVEs(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterCVEHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetCVEHistoryV1(ctx)

	require.NoError(t, err)

	// Collect all CVEs across all entries.
	allCVEs := make(map[string]bool)
	for _, entry := range resp.Entries {
		for _, cve := range entry.CVEs {
			allCVEs[cve] = true
		}
	}

	// The mock HTML contains these specific CVEs.
	assert.True(t, allCVEs["CVE-2026-19012"], "expected CVE-2026-19012 to be found")
	assert.True(t, allCVEs["CVE-2026-19013"], "expected CVE-2026-19013 to be found")
	assert.True(t, allCVEs["CVE-2026-19015"], "expected CVE-2026-19015 to be found")
}
