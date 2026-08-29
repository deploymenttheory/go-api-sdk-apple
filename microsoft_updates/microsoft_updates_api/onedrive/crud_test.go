package onedrive_test

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/onedrive"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/onedrive/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) *onedrive.OneDriveService {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return onedrive.NewService(transport)
}

func TestGetProductionRingV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterProductionManifestMock()

	ctx := context.Background()
	ring, err := svc.GetProductionRingV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, ring)
	assert.Equal(t, onedrive.RingProduction, ring.Ring)
	assert.Equal(t, "26.062.0402", ring.Version)
	assert.NotEmpty(t, ring.DownloadURL)
	assert.Equal(t, onedrive.ApplicationID, ring.ApplicationID)
	assert.Equal(t, onedrive.BundleID, ring.BundleID)
}

func TestGetInsiderRingV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterInsiderManifestMock()

	ctx := context.Background()
	ring, err := svc.GetInsiderRingV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, ring)
	assert.Equal(t, onedrive.RingInsider, ring.Ring)
}

func TestGetProductionRingV1_HTTPError(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterErrorMock("https://g.live.com/0USSDMC_W5T/StandaloneProductManifest")

	ctx := context.Background()
	_, err := svc.GetProductionRingV1(ctx)
	require.Error(t, err)
}

func TestOneDriveConstants(t *testing.T) {
	assert.Equal(t, "ONDR18", onedrive.ApplicationID)
	assert.Equal(t, "com.microsoft.OneDrive", onedrive.BundleID)
}
