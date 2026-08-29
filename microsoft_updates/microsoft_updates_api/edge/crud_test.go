package edge_test

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/edge"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/edge/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T) *edge.EdgeService {
	t.Helper()

	transport, err := client.NewTransport(client.WithRetryCount(0))
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)

	return edge.NewService(transport)
}

func TestGetStableV1_Success(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterStableMock()

	ctx := context.Background()
	resp, err := svc.GetStableV1(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, edge.ChannelStable, resp.Channel)
	require.NotNil(t, resp.Release)
	assert.Equal(t, "147.0.3912.72", resp.Release.Version)
	assert.NotEmpty(t, resp.Release.PublishedTime)
	assert.Len(t, resp.Release.Artifacts, 1)
	assert.NotEmpty(t, resp.Release.Artifacts[0].Location)
	assert.NotEmpty(t, resp.Release.Artifacts[0].HashSHA256)
}

func TestGetStableV1_HTTPError(t *testing.T) {
	svc := setupMockClient(t)
	mocks.RegisterErrorMock(constants.EdgeStableEndpoint)

	ctx := context.Background()
	_, err := svc.GetStableV1(ctx)
	require.Error(t, err)
}

func TestEdgeChannelConstants(t *testing.T) {
	assert.Equal(t, "stable", edge.ChannelStable)
	assert.Equal(t, "beta", edge.ChannelBeta)
	assert.Equal(t, "dev", edge.ChannelDev)
	assert.Equal(t, "canary", edge.ChannelCanary)
}
