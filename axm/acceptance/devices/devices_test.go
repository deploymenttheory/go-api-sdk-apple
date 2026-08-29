package devices

import (
	"context"
	"testing"
	"time"

	acc "github.com/deploymenttheory/go-sdk-appleservices/axm/acceptance"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TestAcceptance_Devices_GetOrganizationDevices
// Verifies the list-all-devices endpoint returns a valid, non-empty response.
// =============================================================================

func TestAcceptance_Devices_GetOrganizationDevices(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	// --- Default options ---
	t.Run("DefaultOptions", func(t *testing.T) {
		acc.LogTestStage(t, "List", "Getting organization devices with default options")

		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.GetV1(ctx1, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Data, "organization should contain at least one device")

		acc.LogTestSuccess(t, "GetOrganizationDevicesV1: found %d devices", len(result.Data))
	})

	// --- Field selection ---
	t.Run("WithFieldSelection", func(t *testing.T) {
		acc.LogTestStage(t, "List", "Getting organization devices with field selection")

		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		opts := &devices.RequestQueryOptions{
			Fields: []string{
				devices.FieldSerialNumber,
				devices.FieldDeviceModel,
				devices.FieldStatus,
				devices.FieldProductFamily,
			},
			Limit: 10,
		}

		result, resp, err := svc.GetV1(ctx2, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		for _, device := range result.Data {
			assert.NotEmpty(t, device.ID, "device ID should not be empty")
			assert.NotNil(t, device.Attributes, "device attributes should not be nil")
		}

		acc.LogTestSuccess(t, "GetOrganizationDevicesV1 (fields): found %d devices", len(result.Data))
	})

	// --- Pagination ---
	// Note: GetOrganizationDevicesV1 fetches all pages automatically, so Limit
	// controls page size sent to the API, not the total result count returned here.
	t.Run("WithPaginationLimit", func(t *testing.T) {
		acc.LogTestStage(t, "List", "Getting organization devices with limit")

		ctx3, cancel3 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel3()

		opts := &devices.RequestQueryOptions{Limit: 2}

		result, resp, err := svc.GetV1(ctx3, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Data, "organization should have at least one device")
		acc.LogTestSuccess(t, "GetOrganizationDevicesV1 (limit=2 page size): %d total device(s) across all pages", len(result.Data))
	})
}

// =============================================================================
// TestAcceptance_Devices_GetDeviceInformationByDeviceID
// Retrieves the first device from the org list and fetches its full record.
// =============================================================================

func TestAcceptance_Devices_GetDeviceInformationByDeviceID(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	// Prerequisite: obtain a real device ID
	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	list, _, err := svc.GetV1(listCtx, &devices.RequestQueryOptions{
		Fields: []string{devices.FieldSerialNumber},
		Limit:  1,
	})
	require.NoError(t, err, "prerequisite: list devices")
	if len(list.Data) == 0 {
		t.Skip("No devices found in organization — skipping device-by-ID acceptance tests")
	}

	deviceID := list.Data[0].ID
	acc.LogTestStage(t, "Setup", "Using device ID=%s for by-ID tests", deviceID)

	// --- All fields ---
	t.Run("AllFields", func(t *testing.T) {
		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.GetByDeviceIDV1(ctx1, deviceID, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.Equal(t, deviceID, result.Data.ID)
		assert.NotNil(t, result.Data.Attributes)

		acc.LogTestSuccess(t, "GetDeviceInformationByDeviceIDV1: ID=%s serial=%s",
			result.Data.ID, result.Data.Attributes.SerialNumber)
	})

	// --- Specific fields ---
	t.Run("SpecificFields", func(t *testing.T) {
		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		opts := &devices.RequestQueryOptions{
			Fields: []string{
				devices.FieldSerialNumber,
				devices.FieldDeviceModel,
				devices.FieldStatus,
			},
		}

		result, resp, err := svc.GetByDeviceIDV1(ctx2, deviceID, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.Equal(t, deviceID, result.Data.ID)

		acc.LogTestSuccess(t, "GetDeviceInformationByDeviceIDV1 (fields): serial=%s model=%s status=%s",
			result.Data.Attributes.SerialNumber,
			result.Data.Attributes.DeviceModel,
			result.Data.Attributes.Status)
	})
}

// =============================================================================
// TestAcceptance_Devices_GetAppleCareInformationByDeviceID
// Verifies AppleCare coverage data is retrievable for a real device.
// =============================================================================

func TestAcceptance_Devices_GetAppleCareInformationByDeviceID(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	// Prerequisite: obtain a real device ID
	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	list, _, err := svc.GetV1(listCtx, &devices.RequestQueryOptions{
		Fields: []string{devices.FieldSerialNumber},
		Limit:  1,
	})
	require.NoError(t, err, "prerequisite: list devices")
	if len(list.Data) == 0 {
		t.Skip("No devices found in organization — skipping AppleCare acceptance tests")
	}

	deviceID := list.Data[0].ID

	t.Run("AllFields", func(t *testing.T) {
		acc.LogTestStage(t, "GetAppleCare", "Getting AppleCare info for device ID=%s", deviceID)

		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		opts := &devices.RequestQueryOptions{
			Fields: []string{
				devices.FieldAppleCareStatus,
				devices.FieldAppleCarePaymentType,
				devices.FieldAppleCareDescription,
				devices.FieldAppleCareStartDateTime,
				devices.FieldAppleCareEndDateTime,
				devices.FieldAppleCareIsRenewable,
				devices.FieldAppleCareIsCanceled,
			},
			Limit: 100,
		}

		result, resp, err := svc.GetAppleCareByDeviceIDV1(ctx1, deviceID, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		acc.LogTestSuccess(t, "GetAppleCareInformationByDeviceIDV1: found %d coverage plan(s) for device %s",
			len(result.Data), deviceID)

		for _, coverage := range result.Data {
			assert.NotEmpty(t, coverage.ID)
			assert.NotNil(t, coverage.Attributes)
			if coverage.Attributes != nil {
				assert.NotEmpty(t, coverage.Attributes.Status)
				assert.NotEmpty(t, coverage.Attributes.Description)
			}
		}
	})

	t.Run("WithNilOptions", func(t *testing.T) {
		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		result, resp, err := svc.GetAppleCareByDeviceIDV1(ctx2, deviceID, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		acc.LogTestSuccess(t, "GetAppleCareInformationByDeviceIDV1 (nil opts): %d plan(s)", len(result.Data))
	})
}

// =============================================================================
// TestAcceptance_Devices_ValidationErrors
// Verifies that client-side validation fires before any HTTP call is made.
// =============================================================================

func TestAcceptance_Devices_ValidationErrors(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	t.Run("GetDeviceInformation_EmptyID", func(t *testing.T) {
		_, _, err := svc.GetByDeviceIDV1(ctx, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "device ID is required")
	})

	t.Run("GetAppleCareInformation_EmptyID", func(t *testing.T) {
		_, _, err := svc.GetAppleCareByDeviceIDV1(ctx, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "device ID is required")
	})
}

// =============================================================================
// TestAcceptance_Devices_ContextCancellation
// Verifies that a cancelled context produces an appropriate error immediately.
// =============================================================================

func TestAcceptance_Devices_ContextCancellation(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.Devices

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // ensure timeout has elapsed

	_, _, err := svc.GetV1(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
