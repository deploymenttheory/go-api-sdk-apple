package appstore_ios_test

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/appstore_ios"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/appstore_ios/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) *appstore_ios.AppStoreIOSService {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return appstore_ios.NewService(transport)
}

func TestGetAppByNameV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterMicrosoftWordMock()

	ctx := context.Background()
	app, err := svc.GetAppByNameV1(ctx, "Microsoft Word")

	require.NoError(t, err)
	require.NotNil(t, app)
	assert.Equal(t, "Microsoft Word", app.TrackName)
	assert.Equal(t, "com.microsoft.Office.Word", app.BundleID)
	assert.NotEmpty(t, app.Version)
}

func TestGetAppByNameV1_EmptyName(t *testing.T) {
	svc := setupMockClient(t)
	ctx := context.Background()

	_, err := svc.GetAppByNameV1(ctx, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestIOSBundleIDConstants(t *testing.T) {
	assert.Equal(t, "com.microsoft.Office.Word", appstore_ios.BundleIDWord)
	assert.Equal(t, "com.microsoft.Office.Excel", appstore_ios.BundleIDExcel)
	assert.Equal(t, "com.microsoft.msedge", appstore_ios.BundleIDEdge)
}
