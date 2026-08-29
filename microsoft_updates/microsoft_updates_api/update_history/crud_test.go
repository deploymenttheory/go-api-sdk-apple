package update_history_test

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/update_history"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/update_history/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) *update_history.UpdateHistoryService {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return update_history.NewService(transport)
}

func TestGetUpdateHistoryV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterUpdateHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetUpdateHistoryV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Entries)

	// At least one entry should have a non-empty release date and version.
	first := resp.Entries[0]
	assert.NotEmpty(t, first.ReleaseDate)
	assert.NotEmpty(t, first.Version)
}

func TestGetUpdateHistoryV1_WithDownloadLinks(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterUpdateHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetUpdateHistoryV1(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Entries)

	// Find an entry with suite download links.
	var fullEntry *update_history.UpdateHistoryEntry
	for i := range resp.Entries {
		if resp.Entries[i].BusinessProSuiteDownload != "" {
			fullEntry = &resp.Entries[i]
			break
		}
	}
	require.NotNil(t, fullEntry, "expected at least one entry with suite download links")
	assert.NotEmpty(t, fullEntry.BusinessProSuiteDownload)
}

func TestGetUpdateHistoryV1_HasPerAppLinks(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterUpdateHistoryMock()

	ctx := context.Background()
	resp, err := svc.GetUpdateHistoryV1(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Entries)

	// Find an entry with per-app update links.
	var fullEntry *update_history.UpdateHistoryEntry
	for i := range resp.Entries {
		if resp.Entries[i].WordUpdate != "" {
			fullEntry = &resp.Entries[i]
			break
		}
	}
	require.NotNil(t, fullEntry, "expected at least one entry with per-app update links")
	assert.NotEmpty(t, fullEntry.WordUpdate)
	assert.NotEmpty(t, fullEntry.ExcelUpdate)
	assert.NotEmpty(t, fullEntry.PowerPointUpdate)
	assert.NotEmpty(t, fullEntry.OutlookUpdate)
	assert.NotEmpty(t, fullEntry.OneNoteUpdate)
	assert.False(t, fullEntry.Archived)
}
