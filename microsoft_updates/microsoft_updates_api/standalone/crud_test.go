package standalone_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/standalone"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/standalone/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) (*standalone.StandaloneService, *client.Transport) {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return standalone.NewService(transport), transport
}

func TestGetPackageByApplicationIDV1_Success(t *testing.T) {
	svc, _ := setupMockClient(t)
	mocks.RegisterWordMock(constants.StandaloneCDNBaseURL)

	ctx := context.Background()
	pkg, err := svc.GetPackageByApplicationIDV1(ctx, standalone.AppIDWord)

	require.NoError(t, err)
	require.NotNil(t, pkg)
	assert.Equal(t, standalone.AppIDWord, pkg.ApplicationID)
	assert.Equal(t, "Microsoft Word", pkg.Title)
	assert.Equal(t, "16.108.1", pkg.ShortVersion)
	assert.Equal(t, "16.108.26041915", pkg.FullVersion)
	assert.Equal(t, "14.0", pkg.MinimumOS)
	assert.NotEmpty(t, pkg.Location)
	assert.NotEmpty(t, pkg.Hash)
	assert.NotEmpty(t, pkg.HashSHA256)
}

func TestGetPackageByApplicationIDV1_EmptyID(t *testing.T) {
	svc, _ := setupMockClient(t)
	ctx := context.Background()

	_, err := svc.GetPackageByApplicationIDV1(ctx, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "application ID is required")
}

func TestGetPackageByApplicationIDV1_HTTPError(t *testing.T) {
	svc, _ := setupMockClient(t)
	mocks.RegisterErrorMock(constants.StandaloneCDNBaseURL + standalone.AppIDWord + ".xml")

	ctx := context.Background()
	_, err := svc.GetPackageByApplicationIDV1(ctx, standalone.AppIDWord)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

func TestGetPackageByNameV1_Success(t *testing.T) {
	svc, _ := setupMockClient(t)
	mocks.RegisterWordMock(constants.StandaloneCDNBaseURL)

	ctx := context.Background()
	pkg, err := svc.GetPackageByNameV1(ctx, "Microsoft Word")

	require.NoError(t, err)
	require.NotNil(t, pkg)
	assert.Equal(t, "Microsoft Word", pkg.Title)
}

func TestGetPackageByNameV1_NotFound(t *testing.T) {
	svc, _ := setupMockClient(t)
	ctx := context.Background()

	_, err := svc.GetPackageByNameV1(ctx, "Nonexistent App")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetPackageByNameV1_EmptyName(t *testing.T) {
	svc, _ := setupMockClient(t)
	ctx := context.Background()

	_, err := svc.GetPackageByNameV1(ctx, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "application name is required")
}

func TestGetLatestV1_PartialSuccess(t *testing.T) {
	svc, _ := setupMockClient(t)

	// Only register Word; all others will 404 and be skipped.
	mocks.RegisterWordMock(constants.StandaloneCDNBaseURL)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(http.StatusNotFound, "not found"))

	ctx := context.Background()
	resp, err := svc.GetLatestV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Packages, 1)
	assert.Equal(t, "Microsoft Word", resp.Packages[0].Title)
}

func TestAppIDConstants(t *testing.T) {
	assert.Equal(t, "MSWD2019", standalone.AppIDWord)
	assert.Equal(t, "XCEL2019", standalone.AppIDExcel)
	assert.Equal(t, "TEAMS21", standalone.AppIDTeams)
}

func TestAppNames(t *testing.T) {
	assert.Equal(t, "Microsoft Word", standalone.AppNames[standalone.AppIDWord])
	assert.Equal(t, "Microsoft Excel", standalone.AppNames[standalone.AppIDExcel])
	assert.Len(t, standalone.AppNames, len(standalone.AllAppIDs))
}
